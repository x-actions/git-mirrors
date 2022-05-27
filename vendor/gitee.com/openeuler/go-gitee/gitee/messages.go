package gitee

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*
Header Example:

Request URL: http://gitee.com/webhook
Request Method: POST
Content-Type: application/json
User-Agent: git-oschina-hook
X-Gitee-Token: ******
X-Gitee-Ping: false
X-Gitee-Event: Merge Request Hook
X-Git-Oschina-Event: Merge Request Hook

*/
const (
	// eventTypeHeader is the Gitee header key used to pass the event type.
	eventTypeHeader = "X-Gitee-Event"
	// tokenTypeHeader is the Gitee header key used to pass the webhook secret.
	tokenTypeHeader = "X-Gitee-Token"
)

var (
	// eventTypeMapping maps webhooks types
	eventTypeMapping = map[string]string{
		"Note Hook":          "NoteEvent",
		"Push Hook":          "PushEvent",
		"Issue Hook":         "IssueEvent",
		"Merge Request Hook": "PullRequestEvent",
		"Tag Push Hook":      "TagPushEvent",
	}
)

// ValidatePayload validates an incoming Gitee Webhook event request
func ValidatePayload(r *http.Request, secretKey []byte) (payload []byte, err error) {
	var body []byte

	switch ct := r.Header.Get("Content-Type"); ct {
	case "application/json":
		var err error
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, err
		}
		payload = body
	case "application/x-www-form-urlencoded":
		const payloadFormParam = "payload"
		var err error
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, err
		}
		form, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		payload = []byte(form.Get(payloadFormParam))
	default:
		return nil, fmt.Errorf("webhook request has unsupported Content-Type %q", ct)
	}

	token := r.Header.Get(tokenTypeHeader)
	if token != string(secretKey) {
		return nil, errors.New("payload token check failed")
	}
	return payload, nil
}

// WebHookType returns the event type of webhook request r.
func WebHookType(r *http.Request) string {
	return r.Header.Get(eventTypeHeader)
}

// ParseWebHook parses the event payload.
func ParseWebHook(messageType string, payload []byte) (interface{}, error) {
	eventType, ok := eventTypeMapping[messageType]
	if !ok {
		return nil, fmt.Errorf("unknown X-Gitee-Event in message: %v", messageType)
	}
	hookEvent := HookEvent{
		Type:       &eventType,
		RawPayload: (*json.RawMessage)(&payload),
	}
	return hookEvent.ParsePayload()
}
