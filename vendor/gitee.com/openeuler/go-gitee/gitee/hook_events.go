package gitee

import (
	"encoding/json"
	"time"
)

// HookEvent represents a Gitee hook event.
type HookEvent struct {
	Type       *string          `json:"type,omitempty"`
	RawPayload *json.RawMessage `json:"payload,omitempty"`
	Actor      *User            `json:"actor,omitempty"`
	CreatedAt  *time.Time       `json:"created_at,omitempty"`
	ID         *string          `json:"id,omitempty"`
}

// ParsePayload parses the event payload. For recognized event types,
// a value of the corresponding struct type will be returned.
func (e *HookEvent) ParsePayload() (payload interface{}, err error) {
	switch *e.Type {
	case "NoteEvent":
		payload = &NoteEvent{}
	case "PushEvent":
		payload = &PushEvent{}
	case "IssueEvent":
		payload = &IssueEvent{}
	case "PullRequestEvent":
		payload = &PullRequestEvent{}
	case "TagPushEvent":
		payload = &TagPushEvent{}
	}

	err = json.Unmarshal(*e.RawPayload, &payload)
	return payload, err
}
