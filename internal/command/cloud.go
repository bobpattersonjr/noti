package command

import (
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/bobpattersonjr/noti/service/bark"
	"github.com/bobpattersonjr/noti/service/bearychat"
	"github.com/bobpattersonjr/noti/service/chanify"
	"github.com/bobpattersonjr/noti/service/gchat"
	"github.com/bobpattersonjr/noti/service/keybase"
	"github.com/bobpattersonjr/noti/service/mattermost"
	"github.com/bobpattersonjr/noti/service/ntfy"
	"github.com/bobpattersonjr/noti/service/pushbullet"
	"github.com/bobpattersonjr/noti/service/pushover"
	"github.com/bobpattersonjr/noti/service/pushsafer"
	"github.com/bobpattersonjr/noti/service/simplepush"
	"github.com/bobpattersonjr/noti/service/slack"
	"github.com/bobpattersonjr/noti/service/telegram"
	"github.com/bobpattersonjr/noti/service/twilio"
	"github.com/bobpattersonjr/noti/service/webhook"
	"github.com/bobpattersonjr/noti/service/zulip"
	"github.com/spf13/viper"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func getBearyChat(title, message string, v *viper.Viper) notification {
	return &bearychat.Notification{
		Text:            fmt.Sprintf("**%s**\n%s", title, message),
		IncomingHookURI: v.GetString("bearychat.incomingHookURI"),
		Client:          httpClient,
	}
}

func getKeybase(title, message string, v *viper.Viper) notification {
	var explodeTime time.Duration
	if v.GetString("keybase.explodingLifetime") != "" {
		// Error handling: if explodingLifetime is set to a unparseable duration,
		// viper will assign it to zero. Replace with -1, which will cause an early
		// error, to ensure the command does not send a regular message on accident.
		// Keybase's exploding messages have stricter security guarantees.
		explodeTime = v.GetDuration("keybase.explodingLifetime")
		if explodeTime == 0 {
			explodeTime = -1
		}
	}

	return &keybase.Notification{
		Conversation:      v.GetString("keybase.conversation"),
		ChannelName:       v.GetString("keybase.channel"),
		Public:            v.GetBool("keybase.public"),
		ExplodingLifetime: explodeTime,
		Message:           fmt.Sprintf("**%s**\n%s", title, message),
	}
}

func getPushbullet(title, message string, v *viper.Viper) notification {
	return &pushbullet.Notification{
		Title:       title,
		Body:        message,
		Type:        "note",
		AccessToken: v.GetString("pushbullet.accessToken"),
		DeviceIden:  v.GetString("pushbullet.deviceIden"),
		Client:      httpClient,
	}
}

func getPushover(title, message string, v *viper.Viper) notification {
	return &pushover.Notification{
		Title:    title,
		Message:  message,
		APIToken: v.GetString("pushover.apiToken"),
		UserKey:  v.GetString("pushover.userKey"),
		Client:   httpClient,
	}
}

func getPushsafer(title, message string, v *viper.Viper) notification {
	return &pushsafer.Notification{
		Title:   title,
		Message: message,
		Key:     v.GetString("pushsafer.key"),
		Client:  httpClient,
	}
}

func getSimplepush(title, message string, v *viper.Viper) notification {
	return &simplepush.Notification{
		Title:   title,
		Message: message,
		Key:     v.GetString("simplepush.key"),
		Event:   v.GetString("simplepush.event"),
		Client:  httpClient,
	}
}

func getSlack(title, message string, v *viper.Viper) notification {
	text := fmt.Sprintf("%s\n%s", title, message)
	if title == v.GetString("slack.username") {
		text = message
	}

	return &slack.Notification{
		Token:     v.GetString("slack.token"),
		Channel:   v.GetString("slack.channel"),
		Username:  v.GetString("slack.username"),
		AppURL:    v.GetString("slack.appurl"),
		Text:      text,
		IconEmoji: ":rocket:",

		Client: httpClient,
	}
}

func getGChat(title, message string, v *viper.Viper) notification {
	return &gchat.Notification{
		Message:  message,
		Title:    title,
		Template: v.GetString("gchat.template"),
		AppURL:   v.GetString("gchat.appurl"),
		Client:   httpClient,
	}
}

func getMattermost(title, message string, v *viper.Viper) notification {
	return &mattermost.Notification{
		IncomingHookURI: v.GetString("mattermost.incomingHookURI"),
		Channel:         v.GetString("mattermost.channel"),
		Username:        v.GetString("mattermost.username"),
		Text:            fmt.Sprintf("**%s %s**\n%s", title, ":rocket:", message),
		IconURL:         v.GetString("mattermost.iconurl"),
		Type:            v.GetString("mattermost.type"),

		Client: httpClient,
	}
}

func getTelegram(title, message string, v *viper.Viper) notification {
	return &telegram.Notification{
		ChatID:  v.GetString("telegram.chatId"),
		Token:   v.GetString("telegram.token"),
		Topic:   v.GetString("telegram.topic"),
		Message: fmt.Sprintf("<b>%s %s</b>\n%s", html.EscapeString(title), "🚀:", message),

		Client: httpClient,
	}
}

func getZulip(title, message string, v *viper.Viper) notification {
	return &zulip.Notification{
		BotAPIKey:       v.GetString("zulip.key"),
		BotEmailAddress: v.GetString("zulip.botAddress"),
		Endpoint:        v.GetString("zulip.URI"),
		Content:         fmt.Sprintf("%s:%s", title, message),
		Type:            v.GetString("zulip.type"),
		To:              v.GetString("zulip.to"),
		Client:          httpClient,
	}
}

func getTwilio(title, message string, v *viper.Viper) notification {
	return &twilio.Notification{
		Content:    fmt.Sprintf("%s:%s", title, message),
		NumberTo:   v.GetString("twilio.numberTo"),
		NumberFrom: v.GetString("twilio.numberFrom"),
		AccountSid: v.GetString("twilio.accountSid"),
		AuthToken:  v.GetString("twilio.authToken"),
	}
}

func getChanify(title, message string, v *viper.Viper) notification {
	return &chanify.Notification{
		ChannelURL:        v.GetString("chanify.channelURL"),
		Text:              message,
		Title:             title,
		Sound:             v.GetBool("chanify.sound"),
		Priority:          v.GetInt("chanify.priority"),
		InterruptionLevel: v.GetString("chanify.interruptionLevel"),
		Client:            httpClient,
	}
}

func getNtfy(title, message string, v *viper.Viper) notification {
	return &ntfy.Notification{
		URL:     v.GetString("ntfy.url"),
		Token:   v.GetString("ntfy.token"),
		Topic:   v.GetString("ntfy.topic"),
		Title:   title,
		Message: message,
		Client:  httpClient,
	}
}

func getBark(title, message string, v *viper.Viper) notification {
	return &bark.Notification{
		URL:       v.GetString("bark.apiurl"),
		DeviceKey: v.GetString("bark.key"),
		Title:     title,
		Body:      message,
		Client:    httpClient,
	}
}

// getWebhooks builds one or more webhook notifications from configuration.
func getWebhooks(title, message string, v *viper.Viper) []notification {
	// Prefer list form under top-level `webhooks:` if provided.
	if raw := v.Get("webhooks"); raw != nil {
		buildFromMap := func(m map[string]interface{}) *webhook.Notification {
			// Headers conversion
			headers := map[string]string{}
			if hv, ok := m["headers"].(map[string]interface{}); ok {
				for k, v := range hv {
					switch vv := v.(type) {
					case string:
						headers[k] = vv
					default:
						headers[k] = fmt.Sprint(vv)
					}
				}
			} else if hs, ok := m["headers"].(map[string]string); ok {
				headers = hs
			}

			// String helpers with fallback to empty
			s := func(key string) string {
				if val, ok := m[key].(string); ok {
					return val
				}
				return ""
			}

			return &webhook.Notification{
				URL:         s("url"),
				Method:      s("method"),
				Headers:     headers,
				ContentType: s("contentType"),
				Template:    s("template"),
				Title:       title,
				Message:     message,
				Client:      httpClient,
			}
		}

		// Case 1: explicitly set as []map[string]interface{}
		if items, ok := raw.([]map[string]interface{}); ok && len(items) > 0 {
			notis := make([]notification, 0, len(items))
			for _, m := range items {
				notis = append(notis, buildFromMap(m))
			}
			return notis
		}
		// Case 2: generic []interface{} slice containing maps
		if items, ok := raw.([]interface{}); ok && len(items) > 0 {
			notis := make([]notification, 0, len(items))
			for _, it := range items {
				m, ok := it.(map[string]interface{})
				if !ok {
					continue
				}
				notis = append(notis, buildFromMap(m))
			}
			if len(notis) > 0 {
				return notis
			}
		}
	}

	// Fallback to single webhook.* configuration.
	return []notification{getWebhook(title, message, v)}
}

func getWebhook(title, message string, v *viper.Viper) notification {
	headers := map[string]string{}
	if m := v.GetStringMapString("webhook.headers"); len(m) > 0 {
		headers = m
	}

	return &webhook.Notification{
		URL:         v.GetString("webhook.url"),
		Method:      v.GetString("webhook.method"),
		Headers:     headers,
		ContentType: v.GetString("webhook.contentType"),
		Template:    v.GetString("webhook.template"),
		Title:       title,
		Message:     message,
		Client:      httpClient,
	}
}
