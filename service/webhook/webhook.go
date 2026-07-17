package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
)

// Notification is a generic Webhook notification.
// It supports configurable method, headers, content type, and a body template.
type Notification struct {
	// URL is the destination webhook endpoint.
	URL string

	// Method is the HTTP method (e.g., POST, PUT, PATCH, GET). Defaults to POST.
	Method string

	// Headers are additional HTTP headers to include in the request.
	Headers map[string]string

	// ContentType sets the Content-Type header when a body is sent.
	ContentType string

	// Template is the body template; it can reference {{.title}} and {{.message}}.
	Template string

	// Title and Message are the notification contents available to the template.
	Title   string
	Message string

	// Client is the HTTP client used to send the request.
	Client *http.Client
}

// Send triggers the webhook notification.
func (n *Notification) Send() error {
	if n.URL == "" {
		return errors.New("webhook: missing URL")
	}

	method := strings.ToUpper(strings.TrimSpace(n.Method))
	if method == "" {
		method = http.MethodPost
	}

	ctype := n.contentType()

	var body io.Reader
	if n.Template != "" && method != http.MethodGet && method != http.MethodHead {
		t, err := template.New("webhook").Parse(n.Template)
		if err != nil {
			return fmt.Errorf("webhook: parse template: %w", err)
		}

		title, message := n.Title, n.Message
		// For JSON payloads, escape the values so a quote or newline in the
		// title or message can't produce an invalid body.
		if isJSON(ctype) {
			title = jsonEscape(title)
			message = jsonEscape(message)
		}

		buf := &bytes.Buffer{}
		if err := t.Execute(buf, map[string]string{
			"title":   title,
			"message": message,
		}); err != nil {
			return fmt.Errorf("webhook: execute template: %w", err)
		}
		body = buf
	}

	req, err := http.NewRequest(method, n.URL, body)
	if err != nil {
		return fmt.Errorf("webhook: new request: %w", err)
	}

	// Apply headers.
	for k, v := range n.Headers {
		if k == "" {
			continue
		}
		req.Header.Set(k, v)
	}

	// Set Content-Type if we have a body and no explicit header already.
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", ctype)
	}

	resp, err := n.Client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook: %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	return nil
}

// contentType resolves the effective Content-Type: an explicit header wins,
// then the ContentType field, then application/json.
func (n *Notification) contentType() string {
	for k, v := range n.Headers {
		if strings.EqualFold(k, "Content-Type") {
			return v
		}
	}
	if ct := strings.TrimSpace(n.ContentType); ct != "" {
		return ct
	}
	return "application/json"
}

func isJSON(contentType string) bool {
	ct := strings.ToLower(contentType)
	return strings.Contains(ct, "application/json") || strings.HasSuffix(strings.SplitN(ct, ";", 2)[0], "+json")
}

// jsonEscape returns s escaped for embedding inside a JSON string literal,
// without the surrounding quotes.
func jsonEscape(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		return s
	}
	return string(b[1 : len(b)-1])
}
