# Boba Bash Shipifyer

Are you a Boba Bash Org? yes.
Are your attendees not shipping, and they wont read your emails? yes.
Are you having a hard time DMing all your attendees personally? no?? (laughs in boba bash delhi has more signups than you)

A simple tool to automatically get your attendees emails, then their slack ids and then DMs them!

## Setup 

### Slack app

create a app with this manifest

```yaml
display_information:
  name: boba-bash-shipifyer
features:
  bot_user:
    display_name: boba-bash-shipifyer
    always_online: false
oauth_config:
  scopes:
    user:
      - im:history
      - chat:write
      - im:write
    bot:
      - users:read.email
      - users:read
  pkce_enabled: false
settings:
  interactivity:
    is_enabled: true
  org_deploy_enabled: false
  socket_mode_enabled: true
  token_rotation_enabled: false
  app_level_token_rotation_enabled: false
  is_mcp_enabled: false
```

and install it, save the bot and user token.

### Org portal

open your organize portal, save the event number (bash.hackclub.com/organize/XX) 

open dev tools, then open the cookies section (via the application tab on chromium and storage tab on firefox) save the "_session_id" cookie

### Message

run 

```bash
mkdir ~/.boba-bash-shipifyer/
nano ~/.boba-bash-shipifyer/message.md
```

and add a message there like "hey, lock the FUCK IN YOU-" /j

### Tool

you need Go and a Keyring setup before hand

install: `go install github.com/reztheperson/boba-bash-shipifyer`

NOTE: IT WILL AUTO SEND MESSAGES TO PARTICIPANTS THAT YOU HAVE NEVER INTERACTED WITH BEFORE, for people you have dms with before it asks.

running: `boba-bash-shipifyer`

on the first time paste the secrets.

enjoy