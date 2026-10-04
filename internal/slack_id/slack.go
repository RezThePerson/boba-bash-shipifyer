package slack_id

import (
	"fmt"
	"log"
	"time"

	"github.com/reztheperson/boba-bash-shipifyer/internal/keyring"
	"github.com/schollz/progressbar/v3"
	"github.com/slack-go/slack"
)

func Get(emails []string) []string {
	api := slack.New(keyring.Secrets.SlackBotToken)

	var slackIDs []string
	var missingEmails []string
	total := len(emails)

	if total == 0 {
		log.Println("[slack] No emails provided to check.")
		return slackIDs
	}

	bar := progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(25),
		progressbar.OptionSetDescription("[cyan]Checking Slack IDs...[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for _, email := range emails {
		// Throttle slightly to respect Slack Tier 3 rate limit (~50 req/min)
		time.Sleep(1200 * time.Millisecond)

		user, err := api.GetUserByEmail(email)
		if err != nil {
			// Catch rate limiting and retry once
			if rateLimitedErr, ok := err.(*slack.RateLimitedError); ok {
				retryAfter := rateLimitedErr.RetryAfter
				if retryAfter == 0 {
					retryAfter = 5 * time.Second
				}
				time.Sleep(retryAfter)
				user, err = api.GetUserByEmail(email)
			}
		}

		if err != nil || user == nil || user.ID == "" {
			missingEmails = append(missingEmails, email)
		} else {
			slackIDs = append(slackIDs, user.ID)
		}

		_ = bar.Add(1)
	}

	fmt.Println() // Newline after progress bar finishes

	log.Printf("[slack] Lookup complete! Found: %d | Missing: %d | Total: %d\n",
		len(slackIDs),
		len(missingEmails),
		total,
	)

	if len(missingEmails) > 0 {
		fmt.Printf("\n--- Emails not found on Slack (%d) ---\n", len(missingEmails))
		for i, email := range missingEmails {
			fmt.Printf("%d. %s\n", i+1, email)
		}
		fmt.Println("---------------------------------------")
	}

	return slackIDs
}