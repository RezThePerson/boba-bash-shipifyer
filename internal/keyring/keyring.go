package keyring

import (
	"errors"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/zalando/go-keyring"
)

var Secrets struct {
	SlackBotToken  string
	SlackUserToken string
	PanelOrgNumber string
	PanelOrgCookie string
}

func Get() {
	Secrets.SlackBotToken = getOrPrompt("slack-bot-token", "Enter Slack Bot Token")
	Secrets.SlackUserToken = getOrPrompt("slack-user-token", "Enter Slack User Token")
	Secrets.PanelOrgNumber = getOrPrompt("panel-org-number", "Enter Panel Org Number (bash.hackclub.com/organize/XX)")
	Secrets.PanelOrgCookie = getOrPrompt("panel-org-cookie", "Enter Panel Org Cookie (XSRF-TOKEN)")
}

func getOrPrompt(key, label string) string {
	// Attempt to get from keyring
	val, err := keyring.Get("boba-bash-shopifyer", key)
	if err == nil {
		return val
	}

	// Attempt to get from prompting user
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("cookie cannot be empty")
			}
			return nil
		},
	}

	val, err = prompt.Run()
	if err != nil {
		log.Fatalf("Error reading %s prompt: %v", label, err)
	}

	// Saving to keyring
	if err = keyring.Set("boba-bash-shopifyer", key, val); err != nil {
		log.Fatalf("Error saving %s to keyring: %v", label, err)
	}

	log.Printf("[keyring] secrets initialized successfully.")

	return val
}
