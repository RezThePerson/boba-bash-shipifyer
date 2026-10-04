package slack

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/reztheperson/boba-bash-shipifyer/internal/keyring"
	"github.com/schollz/progressbar/v3"
	"github.com/slack-go/slack"
)

// getOutboundMessage reads ~/.boba-bash-shipifyer/message.md if present, defaulting to "hi"
func getOutboundMessage() string {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		filePath := filepath.Join(homeDir, ".boba-bash-shipifyer", "message.md")
		if content, err := os.ReadFile(filePath); err == nil {
			trimmed := strings.TrimSpace(string(content))
			if trimmed != "" {
				log.Printf("[messenger] Loaded message template from %s", filePath)
				return trimmed
			}
		}
	}

	log.Println("[messenger] message.md not found or empty; defaulting to 'hi'")
	return "hi"
}

// MessageSlackUsers checks chat history and sends the configured message.
func Msg(slackIDs []string) {
	if len(slackIDs) == 0 {
		log.Println("[messenger] No Slack IDs provided to message.")
		return
	}

	messageText := getOutboundMessage()
	api := slack.New(keyring.Secrets.SlackUserToken)

	total := len(slackIDs)
	sentCount := 0
	skippedCount := 0
	failedCount := 0

	bar := progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(25),
		progressbar.OptionSetDescription("[yellow]Processing DMs...[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for _, userID := range slackIDs {
		time.Sleep(1200 * time.Millisecond)

		// 1. Open or locate direct conversation
		channel, _, _, err := api.OpenConversation(&slack.OpenConversationParameters{
			Users: []string{userID},
		})
		if err != nil {
			log.Printf("[messenger] Could not open DM with %s: %v", userID, err)
			failedCount++
			_ = bar.Add(1)
			continue
		}

		// 2. Fetch history
		historyParams := &slack.GetConversationHistoryParameters{
			ChannelID: channel.ID,
			Limit:     1,
		}

		history, err := api.GetConversationHistory(historyParams)
		if err != nil {
			log.Printf("[messenger] Failed to fetch history for %s: %v", userID, err)
			failedCount++
			_ = bar.Add(1)
			continue
		}

		shouldSend := true

		// 3. Prompt if previous history exists
		if len(history.Messages) > 0 {
			fmt.Print("\r\033[K") // Clear current terminal line for the prompt

			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("User %s has prior conversation history. Send message anyway?", userID),
				IsConfirm: true,
			}

			result, promptErr := prompt.Run()
			if promptErr != nil || strings.ToLower(result) != "y" {
				shouldSend = false
				skippedCount++
				log.Printf("[messenger] Skipped user %s.", userID)
			}
		}

		// 4. Send message
		if shouldSend {
			_, _, postErr := api.PostMessage(
				channel.ID,
				slack.MsgOptionText(messageText, false),
				slack.MsgOptionAsUser(true),
			)
			if postErr != nil {
				log.Printf("[messenger] Failed to send message to %s: %v", userID, postErr)
				failedCount++
			} else {
				sentCount++
			}
		}

		_ = bar.Add(1)
	}

	fmt.Println()
	log.Printf("[messenger] Complete! Sent: %d | Skipped: %d | Failed: %d | Total: %d\n",
		sentCount, skippedCount, failedCount, total)
}
