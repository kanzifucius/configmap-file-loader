package controller

import (
"time"
)

// Config is the controller configuration.
type Config struct {
	ResyncPeriod  time.Duration
	Namespace     string
	Directory     string
	Name          string
	KeyNameFilter string
	LabelFilter   string
	Webhook string
	WebhookMethod string
	WebhookStatusCode int

}
