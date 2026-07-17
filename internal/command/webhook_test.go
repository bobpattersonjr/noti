package command

import (
	"testing"

	wh "github.com/bobpattersonjr/noti/service/webhook"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

func TestGetWebhooksMultiple(t *testing.T) {
	v := viper.New()
	v.Set("webhooks", []map[string]interface{}{
		{
			"url":         "https://example.com/hook1",
			"method":      "POST",
			"contentType": "application/json",
			"template":    "{\"title\":\"{{.title}}\",\"message\":\"{{.message}}\"}",
			"headers": map[string]interface{}{
				"X-Test": "one",
			},
		},
		{
			"url":         "https://example.com/hook2",
			"method":      "PUT",
			"contentType": "text/plain",
			"template":    "{{.title}}: {{.message}}",
		},
	})

	notis := getWebhooks("T", "M", v)
	if len(notis) != 2 {
		t.Fatalf("unexpected webhook notifications count: have=%d want=%d", len(notis), 2)
	}
}

func TestGetNotificationsWithWebhooks(t *testing.T) {
	v := viper.New()
	v.Set("webhooks", []map[string]interface{}{
		{"url": "https://example.com/a"},
		{"url": "https://example.com/b"},
	})

	services := map[string]struct{}{"webhook": {}}
	notis := getNotifications(v, services)
	if len(notis) != 2 {
		t.Fatalf("unexpected notifications count: have=%d want=%d", len(notis), 2)
	}
}

func TestGetNotificationsWebhookSingleFallback(t *testing.T) {
	v := viper.New()
	v.Set("webhook.url", "https://example.com/only")

	services := map[string]struct{}{"webhook": {}}
	notis := getNotifications(v, services)
	if len(notis) != 1 {
		t.Fatalf("unexpected notifications count: have=%d want=%d", len(notis), 1)
	}
}

func TestWebhookHeadersFromEnv(t *testing.T) {
	// Save and restore env
	origHeaders := os.Getenv("NOTI_WEBHOOK_HEADERS")
	origURL := os.Getenv("NOTI_WEBHOOK_URL")
	defer func() {
		os.Setenv("NOTI_WEBHOOK_HEADERS", origHeaders)
		os.Setenv("NOTI_WEBHOOK_URL", origURL)
	}()

	os.Setenv("NOTI_WEBHOOK_HEADERS", "X-A=1, X-B=two")
	os.Setenv("NOTI_WEBHOOK_URL", "https://example.com/env")

	v := viper.New()
	flags := pflag.NewFlagSet("envwebhook", pflag.ContinueOnError)
	InitFlags(flags)
	if err := configureApp(v, flags); err != nil {
		t.Fatalf("configureApp error: %v", err)
	}

	notis := getWebhooks("T", "M", v)
	if len(notis) != 1 {
		t.Fatalf("unexpected notifications count: have=%d want=%d", len(notis), 1)
	}

	// Type assert to webhook.Notification to inspect headers
	n, ok := notis[0].(*wh.Notification)
	if !ok {
		t.Fatalf("unexpected type: %T", notis[0])
	}
	if n.Headers["X-A"] != "1" || n.Headers["X-B"] != "two" {
		t.Fatalf("unexpected headers: %#v", n.Headers)
	}
}
