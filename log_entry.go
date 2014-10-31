package txlog

import (
	"strings"
	"time"
)

type logEntry struct {
	Tag       string    `json:"tag,omitempty"`
	Hostname  string    `json:"host"`
	Service   string    `json:"service"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Context   Context   `json:"context,omitempty"`
}

func NewLogEntry(hostname, service, message string) logEntry {
	return logEntry{
		Hostname:  hostname,
		Service:   service,
		Message:   strings.TrimSpace(message),
		Timestamp: time.Now(),
	}
}
