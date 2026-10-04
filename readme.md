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


