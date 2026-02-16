package store

import "time"

type AuditEntry struct {
	Event     string
	CreatedAt time.Time
}

func NewAuditEntry(event string) AuditEntry {
	return AuditEntry{Event: event, CreatedAt: time.Now().UTC()}
}
