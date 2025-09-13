package webhook

import (
    "bytes"
    "errors"
    "fmt"
    "html/template"
    "io"
    "net/http"
    "strings"
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

    var body io.Reader
    if n.Template != "" && method != http.MethodGet && method != http.MethodHead {
        t, err := template.New("webhook").Parse(n.Template)
        if err != nil {
            return fmt.Errorf("webhook: parse template: %w", err)
        }
        buf := &bytes.Buffer{}
        if err := t.Execute(buf, map[string]string{
            "title":   n.Title,
            "message": n.Message,
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
        ctype := n.ContentType
        if strings.TrimSpace(ctype) == "" {
            ctype = "application/json"
        }
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

