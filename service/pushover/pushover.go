package pushover

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"
)

var (
    // API is the Pushover API endpoint.
    API = "https://api.pushover.net/1/messages.json"
)

type apiResponse struct {
    Info    string   `json:"info"`
    Status  int      `json:"status"`
    Request string   `json:"request"`
    Errors  []string `json:"errors"`
    Token   string   `json:"token"`
}

// Notification is a pushover notification.
type Notification struct {
    Message  string
    Title    string
    APIToken string
    UserKey  string

    Sound     string
    Device    string
    Priority  int // valid values: -2..2
    URL       string
    URLTitle  string
    HTML      int // 0 or 1
    Timestamp int64

    Client *http.Client `json:"-"`
}

// Send sends a pushover notification.
func (n *Notification) Send() error {
    if n.Client == nil {
        n.Client = http.DefaultClient
    }

    vals := make(url.Values)
    vals.Set("token", n.APIToken)
    vals.Set("user", n.UserKey)
    vals.Set("message", n.Message)
    if n.Title != "" {
        vals.Set("title", n.Title)
    }
    if n.Sound != "" {
        vals.Set("sound", n.Sound)
    }
    if n.Device != "" {
        vals.Set("device", n.Device)
    }
    if n.Priority != 0 {
        vals.Set("priority", strconv.Itoa(n.Priority))
    }
    if n.URL != "" {
        vals.Set("url", n.URL)
    }
    if n.URLTitle != "" {
        vals.Set("url_title", n.URLTitle)
    }
    if n.HTML != 0 {
        vals.Set("html", strconv.Itoa(n.HTML))
    }
    // Add a timestamp to avoid out-of-order messages unless already set.
    ts := n.Timestamp
    if ts == 0 {
        ts = time.Now().Unix()
    }
    vals.Set("timestamp", strconv.FormatInt(ts, 10))

    resp, err := n.Client.PostForm(API, vals)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var r apiResponse
    if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
        return err
    }

    if r.Status != 1 {
        return errors.New(strings.Join(r.Errors, ": "))
    } else if strings.Contains(r.Info, "no active devices") {
        return errors.New(r.Info)
    }

    return nil
}
