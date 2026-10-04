package main

import (
	"log"

	"github.com/reztheperson/boba-bash-shipifyer/internal/emails"
	"github.com/reztheperson/boba-bash-shipifyer/internal/keyring"
	"github.com/reztheperson/boba-bash-shipifyer/internal/slack_id"
)

func main() {
	keyring.Get()

	emails := emails.Fetch()

	slackIDs := slack_id.Get(emails)
	log.Printf("Resolved %d Slack user IDs: %v", len(slackIDs), slackIDs)
}
