package pushover

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	n := Notification{
		Title:    "title",
		Message:  "mesg",
		APIToken: "tok",
		UserKey:  "dst",
		Client:   &http.Client{Timeout: 3 * time.Second},
	}
	var mockResp apiResponse
	var hitServer bool

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hitServer = true

		if r.Method != "POST" {
			t.Error("HTTP method should be POST")
		}

		if r.FormValue("token") == "" {
			t.Error("missing access token")
		}
		if r.FormValue("user") == "" {
			t.Error("missing destination")
		}

		json.NewEncoder(rw).Encode(mockResp)
	}))
	defer ts.Close()

	API = ts.URL
	mockResp.Status = 1 // success
	if err := n.Send(); err != nil {
		t.Error(err)
	}

	if !hitServer {
		t.Error("didn't reach server")
	}

	mockResp.Status = 0 // failure
	mockResp.Errors = []string{"error"}
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}

	mockResp.Status = 1 // failure
	mockResp.Info = "no active devices to send to"
	if err := n.Send(); err == nil {
		t.Error("unexpected success")
	}
}

func TestSendExtras(t *testing.T) {
	n := Notification{
		Title:    "title",
		Message:  "mesg",
		APIToken: "tok",
		UserKey:  "dst",
		Sound:    "cosmic",
		Device:   "iphone",
		Priority: 2,
		Retry:    30,
		Expire:   3600,
		Client:   &http.Client{Timeout: 3 * time.Second},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		wantForm := map[string]string{
			"sound":    "cosmic",
			"device":   "iphone",
			"priority": "2",
			"retry":    "30",
			"expire":   "3600",
		}
		for k, want := range wantForm {
			if have := r.FormValue(k); have != want {
				t.Errorf("form value %s: have=%q; want=%q", k, have, want)
			}
		}

		json.NewEncoder(rw).Encode(apiResponse{Status: 1})
	}))
	defer ts.Close()

	API = ts.URL
	if err := n.Send(); err != nil {
		t.Error(err)
	}
}

func TestSendExtrasOmittedWhenUnset(t *testing.T) {
	n := Notification{
		Title:    "title",
		Message:  "mesg",
		APIToken: "tok",
		UserKey:  "dst",
		Client:   &http.Client{Timeout: 3 * time.Second},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for _, k := range []string{"sound", "device", "priority", "retry", "expire"} {
			if _, ok := r.PostForm[k]; ok {
				t.Errorf("form value %s should not be set", k)
			}
		}

		json.NewEncoder(rw).Encode(apiResponse{Status: 1})
	}))
	defer ts.Close()

	API = ts.URL
	if err := n.Send(); err != nil {
		t.Error(err)
	}
}

func TestSendEmergencyPriorityRequiresRetryExpire(t *testing.T) {
	n := Notification{
		Title:    "title",
		Message:  "mesg",
		APIToken: "tok",
		UserKey:  "dst",
		Priority: 2,
		Client:   &http.Client{Timeout: 3 * time.Second},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		t.Error("server should not be reached")
	}))
	defer ts.Close()

	API = ts.URL
	if err := n.Send(); err == nil {
		t.Error("unexpected success: priority 2 without retry and expire")
	}
}
