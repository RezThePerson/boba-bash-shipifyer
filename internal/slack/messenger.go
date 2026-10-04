package slack

import (
	"fmt"
	"log"
	"time"

	"github.com/reztheperson/boba-bash-shipifyer/internal/keyring"
	"github.com/schollz/progressbar/v3"
	"github.com/slack-go/slack"
)

func Msg(slackIDs []string) {
	if len(slackIDs) == 0 {
		log.Println("[messenger] No Slack IDs provided to message.")
		return
	}

	// Initialize SDK with User Token (xoxp-...)
	api := slack.New(keyring.Secrets.SlackUserToken)

	total := len(slackIDs)
	sentCount := 0
	failedCount := 0

	bar := progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(25),
		progressbar.OptionSetDescription("[yellow]Sending DMs...[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for _, userID := range slackIDs {
		// Slack chat.postMessage rate limit is ~1 req/sec per channel/DM
		time.Sleep(1200 * time.Millisecond)

		// Passing userID directly as the channel parameter opens/posts to the DM
		_, _, err := api.PostMessage(
			userID,
			slack.MsgOptionText("hi from the bot, lock in and ship", false),
			slack.MsgOptionAsUser(true), // Send as the authenticated user
		)

		if err != nil {
			// Catch rate limiting and retry once
			if rateLimitedErr, ok := err.(*slack.RateLimitedError); ok {
				retryAfter := rateLimitedErr.RetryAfter
				if retryAfter == 0 {
					retryAfter = 5 * time.Second
				}
				time.Sleep(retryAfter)
				_, _, err = api.PostMessage(userID, slack.MsgOptionText("hi", false), slack.MsgOptionAsUser(true))
			}
		}

		if err != nil {
			log.Printf("[messenger] Failed to message %s: %v", userID, err)
			failedCount++
		} else {
			sentCount++
		}

		_ = bar.Add(1)
	}

	fmt.Println()
	log.Printf("[messenger] Complete! Sent: %d | Failed: %d | Total: %d\n", sentCount, failedCount, total)
}