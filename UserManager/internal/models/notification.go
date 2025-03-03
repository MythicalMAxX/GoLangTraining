package models

type NotificationPriority string

const (
    High   NotificationPriority = "high"
    Medium NotificationPriority = "medium"
    Low    NotificationPriority = "low"
)

type Notification struct {
    Service  string              `json:"service"`
    Message  string              `json:"message"`
    Priority NotificationPriority `json:"priority"`
}