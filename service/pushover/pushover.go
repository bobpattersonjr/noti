package pushover

import (
  "encoding/json"
  "errors"
  "net/http"
  "net/url"
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
  Message  string `json:"message"`
  Title    string `json:"title"`
  APIToken string `json:"token"`
  UserKey  string `json:"user"`
  []Sound    string `json:"sound"`
  []Device   string `json:"device"`
  // Priority is a struct of [-2, -1, 0, 1, 2]
  []Priority int  `json:"priority"`
  []TTL      int  `json:"ttl"`
  []URL      string `json:"url"`
  []URLTitle string `json:"url_title"`
  []HTML     int `json:"html"`
  
  // Send the current time as a timestamp to avoid out of order messages
  []Timestamp int64 `json:"timestamp"`

  Client *http.Client `json:"-"`
}

// Send sends a pushover notification.
func (n *Notification) Send() error {
  vals := make(url.Values)
  vals.Set("token", n.APIToken)
  vals.Set("user", n.UserKey)
  vals.Set("message", n.Message)
  vals.Set("title", n.Title)
  vals.Set("sound", n.Sound)
  vals.Set("device", n.Device)
  vals.Set("priority", n.Priority)

  // Send the current time as a timestamp to avoid out of order messages
  n.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)

  payload, err := json.Marshal(n)
	if err != nil {
		return err
	}

  resp, err := n.Client.PostForm(API, bytes.NewReader(payload))
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
