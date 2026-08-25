package model

import (
	"strings"
	"time"
)

type Bucket struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	RuleID    string    `json:"rule_id"`
	Tokens    float64   `json:"tokens"`
	LastRefill time.Time `json:"last_refill"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Bucket) Validate() error {
	b.ClientID = strings.TrimSpace(b.ClientID)
	b.RuleID = strings.TrimSpace(b.RuleID)
	if b.ClientID == "" {
		return NewValidationError("client_id", "ClientID 不能为空")
	}
	if b.RuleID == "" {
		return NewValidationError("rule_id", "RuleID 不能为空")
	}
	if b.Tokens < 0 {
		return NewValidationError("tokens", "令牌数不能为负数")
	}
	return nil
}

type BucketFilter struct {
	ClientID string
	RuleID   string
}

func (f BucketFilter) Match(b *Bucket) bool {
	if f.ClientID != "" && b.ClientID != f.ClientID {
		return false
	}
	if f.RuleID != "" && b.RuleID != f.RuleID {
		return false
	}
	return true
}
