% NOTI.YAML(5) noti 3.8.2 | Noti Configuration File Format
% variadico
% 2018/03/25

#  NAME

noti.yaml - noti configuration file

# SYNOPSIS

noti.yaml

# DESCRIPTION

File format is YAML.

If not explicitly set with \--file, then noti will check the following paths,
in the following order.

* ./.noti.yaml
* $XDG_CONFIG_HOME/noti/noti.yaml

If $XDG_CONFIG_HOME is empty, then $HOME/.config will be used as its default
value and noti will check $HOME/.config/noti/noti.yaml.

Settings are resolved in the following order, from highest precedence to
lowest.

1. Command line flags
2. This configuration file
3. Environment variables
4. Built-in defaults

A key set in this file overrides the corresponding `NOTI_*` environment
variable. Precedence applies per key: environment variables still take effect
for keys this file does not set. Command line flags override both.

A complete annotated example covering every service is provided in
docs/noti.example.yaml in the noti source distribution.

# BANNER

icon
: Path to notification icon image. On Linux, accepts an image path or a
  freedesktop icon theme name. Not supported on macOS. On Windows, accepts
  an .ico file path.

# NSUSER

soundName
: Banner success sound. Default is Ping. Possible options are Basso, Blow,
  Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
  Submarine, Tink. See /System/Library/Sounds for available sounds.

soundNameFail
: Banner failure sound. Default is Basso. Possible options are Basso,
  Blow, Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
  Submarine, Tink. See /System/Library/Sounds for available sounds.

# SAY

voice
: Name of voice used for speech notifications.

# ESPEAK

voiceName
: Name of voice used for speech notifications.

# SPEECHSYNTHESIZER

voice
: Name of voice used for speech notifications.

# BLINK1

path
: Path to the blink1-tool executable, or a bare name found on PATH.
  Default is blink1-tool.

color
: LED color as an RRGGBB hex string. Default is FF0000.

brightness
: Brightness from 0 to 255. 0 uses the device default.

delay
: Milliseconds to wait between blinks.

fade
: Milliseconds to fade between colors.

glimmer
: Glimmer the LED the given number of times.

random
: Use a random color on each blink. Set to 1 to enable.

repeats
: Number of times to repeat the blink pattern. Default is 3.

# BEARYCHAT

incomingHookURI
: BearyChat incoming URI.

# KEYBASE

conversation
: Keybase message destination. Can be either users (comma-separated) or team.

channel
: Keybase team's chat channel to send to. Conversation must be a team.
  If empty, the team's default channel will be used (typically "general").

explodingLifetime
: Keybase self-destructing message, after the specified time. Times are
  written like `30s` (30 seconds), `15m` (15 minutes), `24h` (24 hours).

public
: Enables broadcasting a message to everyone (when `conversation` is
  your username), or to teams (when `conversation` is your team name).

# PUSHBULLET

accessToken
: Pushbullet access token. Log into your Pushbullet account and retrieve a
  token from the Account Settings page.

deviceIden
: Pushbullet device iden of the target device, if sending to a single device.

# PUSHOVER

apiToken
: Pushover access token. Log into your Pushover account and create a
  token from the Create New Application/Plugin page.

userKey
: Pushover message destination. Should be your User Key.

# PUSHSAFER

key
: Pushsafer private or alias key. Log into your Pushsafer account and note
  your private or alias key.

# SIMPLEPUSH

key
: Simplepush key. Install the Simplepush app and retrieve your key there.

event
: Customize ringtone and vibration.

# SLACK

token
: Slack access token. Log into your Slack account and retrieve a token
  from the Slack Web API page.

channel
: Slack message destination. Can be either a #channel or a @username.

username
: Noti bot username.

# TWILIO

AuthToken
: Twilio access token. Log into your Twilio account and copy the AuthToken from your project dashboard

accountSid
: Twilio account id. Log into your Twilio account and copy the accountSid from your project dashboard.

numberTo
: This parameter determines the destination phone number for your SMS message. Format this number with a '+' and a country code, e.g., +16175551212

numberFrom
: From specifies the Twilio phone number, short code, or Messaging Service that sends this message. This must be a Twilio phone number that you own, formatted with a '+' and country code, e.g. +16175551212 (E.164 format)


# GCHAT

# WEBHOOK

url
: Destination webhook endpoint URL.

method
: HTTP method to use. Defaults to `POST`.

contentType
: Content-Type for the request when a body is sent. Defaults to `application/json`.

template
: Body template rendered with `title` and `message` fields, e.g.
  `{"title":"{{.title}}","message":"{{.message}}"}`.

headers
: Map of additional headers to include in the request.

webhooks
: A list of webhook objects supporting the same fields as above. When present,
  each item will be sent when the webhook service is enabled.

appurl
: This parameter defines the URL for the Google Chat webhook.

template
: This parameter defines the template combining the title and the message. The default is: '*{{.title}}*: {{.message}}'


# EXAMPLES

    ---
    banner:
      icon: /path/to/icon.png
    nsuser:
      soundName: Ping
      soundNameFail: Basso
    say:
      voice: Alex
    espeak:
      voiceName: english-us
    speechsynthesizer:
      voice: Microsoft David Desktop
    blink1:
      color: 00FF00
      repeats: 3
    bearychat:
      incomingHookURI: 1234567890abcdefg
    keybase:
      conversation: yourteam
      channel: general
    pushbullet:
      accessToken: 1234567890abcdefg
      deviceIden: 1234567890abcdefg
    pushover:
      userKey: 1234567890abcdefg
      apiToken: 1234567890abcdefg
    pushsafer:
      key: 1234567890abcdefg
    simplepush:
      key: 1234567890abcdefg
      event: 1234567890abcdefg
    slack:
      appurl: 'https://hooks.slack.com/services/xxx/yyy/zzz'
    webhook:
      url: 'https://example.com/webhook'
      method: POST
      contentType: application/json
      template: '{"title":"{{.title}}","message":"{{.message}}"}'
      headers:
        X-Custom-Header: abc123
    webhooks:
      - url: 'https://hooks.example.com/one'
        method: POST
        contentType: application/json
        template: '{"t":"{{.title}}","m":"{{.message}}"}'
        headers:
          Authorization: Bearer abc123
      - url: 'https://hooks.example.com/two'
        method: PUT
        contentType: text/plain
        template: '{{.title}}: {{.message}}'
    twilio:
      numberto: +972542877978
      numberfrom: +18111119711
      accountsid: AC3cd135aa82XXXXXXXXf792ba23fc98
      authtoken: 74efd0bXXXXXXXXXXX32f7daca
	gchat:
	  appurl: 'https://chat.googleapis.com/v1/spaces/example/messages?key=keyexample'
	  template: '*{{.title}}*: {{.message}}'



# SEE ALSO

noti(1)
