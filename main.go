package main

import (
	"log"

	"github.com/reztheperson/boba-bash-shipifyer/internal/emails"
	"github.com/reztheperson/boba-bash-shipifyer/internal/keyring"
)

func main() {
	keyring.Get()

	emails := emails.Fetch()

	
}
