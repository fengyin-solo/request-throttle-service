package model

import (
	"strings"
	"time"
)

const (
	AlertLevelWarn     = "warn"
	AlertLevelCritical = "critical"
)

type Alert struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	RuleID      string    `json:"rule_id"`
	Level       string    `json:"level"`
	Message     string    `json:"message"`
	TriggeredAt time.Time `json:"triggered_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (a *Alert) Validate() error {
	a.ClientID = strings.TrimSpace(a.ClientID)
	a.RuleID = strings.TrimSpace(a.RuleID)
	a.Message = strings.TrimSpace(a.Message)
	if a.ClientID == "" {
		return NewValidationError("client_id", "ClientID 不能为空")
	}
	if a.RuleID == "" {
		return NewValidationError("rule_id", "RuleID 不能为空")
	}
	if a.Message == "" {
		return NewValidationError("message", "告警消息不能为空")
	}
	if a.Level == "" {
		a.Level = AlertLevelWarn
	}
	if a.Level != AlertLevelWarn && a.Level != AlertLevelCritical {
		return NewValidationError("level", "告警级别不合法")
	}
	return nil
}

type AlertFilter struct {
	ClientID string
	Level    string
}

func (f AlertFilter) Match(a *Alert) bool {
	if f.ClientID != "" && a.ClientID != f.ClientID {
		return false
	}
	if f.Level != "" && a.Level != f.Level {
		return false
	}
	return true
}
