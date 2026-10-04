package main

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/reztheperson/boba-bash-shipifyer/internal/emails"
	"github.com/reztheperson/boba-bash-shipifyer/internal/keyring"
	"github.com/reztheperson/boba-bash-shipifyer/internal/slack"
	"github.com/reztheperson/boba-bash-shipifyer/internal/slack_id"
)

func main() {
	keyring.Get()

	prompt := promptui.Select{
		Label: "What do you want to do?",
		Items: []string{
			"Full flow (get emails, slack id, and send msg)",
			"Get emails + not shipped count (aka test for correct panel cookie)",
			"Get slack id + send msg to custom email (aka test for correct slack tokens)",
		},
	}

	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Aborted")
		os.Exit(1)
	}

	switch idx {
	case 0:
		emails := emails.Fetch()
		slackIDs := slack_id.Get(emails)
		slack.Msg(slackIDs)
	case 1:
		fmt.Println("Fetching Slack IDs...")
	case 2:
		input := promptui.Prompt{
			Label: "Enter custom email to get slack id and send msg to (run /se info @user if you dont know)",
		}
		email, err := input.Run()
		if err != nil {
			fmt.Println("Aborted")
			os.Exit(1)
		}

		emails := []string{email}

		slackIDs := slack_id.Get(emails)
		slack.Msg(slackIDs)
	}
}
