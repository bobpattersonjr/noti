package webhook

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const defaultTemplate = `{"title":"{{.title}}","message":"{{.message}}"}`

type recorded struct {
	method      string
	contentType string
	headers     http.Header
	body        string
}

func serve(t *testing.T, status int) (*httptest.Server, *recorded) {
	t.Helper()

	rec := &recorded{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		rec.method = r.Method
		rec.contentType = r.Header.Get("Content-Type")
		rec.headers = r.Header.Clone()
		rec.body = string(b)
		w.WriteHeader(status)
	}))
	t.Cleanup(ts.Close)
	return ts, rec
}

func TestSend(t *testing.T) {
	ts, rec := serve(t, http.StatusOK)

	n := Notification{
		URL:      ts.URL,
		Template: defaultTemplate,
		Title:    "title",
		Message:  "mess",
		Client:   ts.Client(),
	}
	if err := n.Send(); err != nil {
		t.Fatal(err)
	}

	if rec.method != http.MethodPost {
		t.Errorf("method: have=%s; want=POST", rec.method)
	}
	if rec.contentType != "application/json" {
		t.Errorf("content type: have=%s; want=application/json", rec.contentType)
	}
	if want := `{"title":"title","message":"mess"}`; rec.body != want {
		t.Errorf("body: have=%s; want=%s", rec.body, want)
	}
}

func TestSendJSONEscaping(t *testing.T) {
	ts, rec := serve(t, http.StatusOK)

	title := `build "core" done`
	message := "warnings & 2 < 3\nsecond line\ttabbed \\ backslash"

	n := Notification{
		URL:      ts.URL,
		Template: defaultTemplate,
		Title:    title,
		Message:  message,
		Client:   ts.Client(),
	}
	if err := n.Send(); err != nil {
		t.Fatal(err)
	}

	var got map[string]string
	if err := json.Unmarshal([]byte(rec.body), &got); err != nil {
		t.Fatalf("body is not valid JSON: %s\nbody: %s", err, rec.body)
	}
	if got["title"] != title {
		t.Errorf("title round-trip: have=%q; want=%q", got["title"], title)
	}
	if got["message"] != message {
		t.Errorf("message round-trip: have=%q; want=%q", got["message"], message)
	}
	if strings.Contains(rec.body, "&#34;") || strings.Contains(rec.body, "&amp;") {
		t.Errorf("body contains HTML entities: %s", rec.body)
	}
}

func TestSendNonJSONNotEscaped(t *testing.T) {
	ts, rec := serve(t, http.StatusOK)

	n := Notification{
		URL:         ts.URL,
		ContentType: "text/plain",
		Template:    "{{.title}}: {{.message}}",
		Title:       `a "quoted" title`,
		Message:     "line1\nline2",
		Client:      ts.Client(),
	}
	if err := n.Send(); err != nil {
		t.Fatal(err)
	}

	if want := "a \"quoted\" title: line1\nline2"; rec.body != want {
		t.Errorf("body: have=%q; want=%q", rec.body, want)
	}
	if rec.contentType != "text/plain" {
		t.Errorf("content type: have=%s; want=text/plain", rec.contentType)
	}
}

func TestSendHeadersAndMethod(t *testing.T) {
	ts, rec := serve(t, http.StatusOK)

	n := Notification{
		URL:      ts.URL,
		Method:   "put",
		Template: defaultTemplate,
		Headers: map[string]string{
			"Authorization": "Bearer tok",
			"X-Custom":      "yes",
		},
		Title:   "t",
		Message: "m",
		Client:  ts.Client(),
	}
	if err := n.Send(); err != nil {
		t.Fatal(err)
	}

	if rec.method != http.MethodPut {
		t.Errorf("method: have=%s; want=PUT", rec.method)
	}
	if have := rec.headers.Get("Authorization"); have != "Bearer tok" {
		t.Errorf("auth header: have=%s; want=Bearer tok", have)
	}
	if have := rec.headers.Get("X-Custom"); have != "yes" {
		t.Errorf("custom header: have=%s; want=yes", have)
	}
}

func TestSendExplicitContentTypeHeaderWins(t *testing.T) {
	ts, rec := serve(t, http.StatusOK)

	n := Notification{
		URL:         ts.URL,
		ContentType: "application/json",
		Template:    "{{.title}}",
		Headers:     map[string]string{"Content-Type": "text/plain"},
		Title:       `no "escaping" here`,
		Message:     "m",
		Client:      ts.Client(),
	}
	if err := n.Send(); err != nil {
		t.Fatal(err)
	}

	if rec.contentType != "text/plain" {
		t.Errorf("content type: have=%s; want=text/plain", rec.contentType)
	}
	if want := `no "escaping" here`; rec.body != want {
		t.Errorf("body: have=%q; want=%q", rec.body, want)
	}
}

func TestSendGetHasNoBody(t *testing.T) {
	ts, rec := serve(t, http.StatusOK)

	n := Notification{
		URL:      ts.URL,
		Method:   "GET",
		Template: defaultTemplate,
		Title:    "t",
		Message:  "m",
		Client:   ts.Client(),
	}
	if err := n.Send(); err != nil {
		t.Fatal(err)
	}

	if rec.method != http.MethodGet {
		t.Errorf("method: have=%s; want=GET", rec.method)
	}
	if rec.body != "" {
		t.Errorf("body should be empty for GET; have=%q", rec.body)
	}
	if rec.contentType != "" {
		t.Errorf("content type should be unset for GET; have=%s", rec.contentType)
	}
}

func TestSendNon2xxIsError(t *testing.T) {
	ts, _ := serve(t, http.StatusBadRequest)

	n := Notification{
		URL:      ts.URL,
		Template: defaultTemplate,
		Title:    "t",
		Message:  "m",
		Client:   ts.Client(),
	}
	if err := n.Send(); err == nil {
		t.Error("want error on 400 response; have nil")
	}
}

func TestSendMissingURL(t *testing.T) {
	n := Notification{Template: defaultTemplate}
	if err := n.Send(); err == nil {
		t.Error("want error on missing URL; have nil")
	}
}
