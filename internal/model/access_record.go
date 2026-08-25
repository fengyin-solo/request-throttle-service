package model

import (
	"strings"
	"time"
)

type AccessRecord struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	RuleID    string    `json:"rule_id"`
	Path      string    `json:"path"`
	Allowed   bool      `json:"allowed"`
	Remaining int       `json:"remaining"`
	LatencyMs int       `json:"latency_ms"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *AccessRecord) Validate() error {
	a.ClientID = strings.TrimSpace(a.ClientID)
	a.RuleID = strings.TrimSpace(a.RuleID)
	a.Path = strings.TrimSpace(a.Path)
	if a.ClientID == "" {
		return NewValidationError("client_id", "ClientID 不能为空")
	}
	if a.RuleID == "" {
		return NewValidationError("rule_id", "RuleID 不能为空")
	}
	if a.Path == "" {
		return NewValidationError("path", "Path 不能为空")
	}
	return nil
}

type AccessRecordFilter struct {
	ClientID string
	RuleID   string
	Allowed  *bool
	Path     string
}

func (f AccessRecordFilter) Match(a *AccessRecord) bool {
	if f.ClientID != "" && a.ClientID != f.ClientID {
		return false
	}
	if f.RuleID != "" && a.RuleID != f.RuleID {
		return false
	}
	if f.Allowed != nil && a.Allowed != *f.Allowed {
		return false
	}
	if f.Path != "" && a.Path != f.Path {
		return false
	}
	return true
}
