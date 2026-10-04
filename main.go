package main

import (
	"github.com/reztheperson/boba-bash-shipifyer/internal/emails"
	"github.com/reztheperson/boba-bash-shipifyer/internal/keyring"
	"github.com/reztheperson/boba-bash-shipifyer/internal/slack"
	"github.com/reztheperson/boba-bash-shipifyer/internal/slack_id"
)

func main() {
	keyring.Get()

	emails := emails.Fetch()

	// emails := []string{"reztheperson@proton.me"}

	slackIDs := slack_id.Get(emails)

	slack.Msg(slackIDs)
}
