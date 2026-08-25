package model

import (
	"strings"
	"time"
)

const (
	RuleAlgorithmFixedWindow  = "fixed_window"
	RuleAlgorithmSlidingWindow = "sliding_window"
	RuleAlgorithmTokenBucket   = "token_bucket"
)

const (
	RuleStatusActive   = "active"
	RuleStatusInactive = "inactive"
)

type Rule struct {
	ID         string    `json:"id"`
	Key        string    `json:"key"`
	Algorithm  string    `json:"algorithm"`
	Limit      int       `json:"limit"`
	WindowSec  int       `json:"window_sec"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (r *Rule) Validate() error {
	r.Key = strings.TrimSpace(r.Key)
	if r.Key == "" {
		return NewValidationError("key", "规则 Key 不能为空")
	}
	if r.Algorithm == "" {
		r.Algorithm = RuleAlgorithmFixedWindow
	}
	if r.Algorithm != RuleAlgorithmFixedWindow && r.Algorithm != RuleAlgorithmSlidingWindow && r.Algorithm != RuleAlgorithmTokenBucket {
		return NewValidationError("algorithm", "算法类型不合法")
	}
	if r.Limit <= 0 {
		return NewValidationError("limit", "限流阈值必须大于 0")
	}
	if r.WindowSec <= 0 {
		return NewValidationError("window_sec", "窗口时间必须大于 0")
	}
	if r.Status == "" {
		r.Status = RuleStatusActive
	}
	if r.Status != RuleStatusActive && r.Status != RuleStatusInactive {
		return NewValidationError("status", "规则状态不合法")
	}
	return nil
}

type RuleFilter struct {
	Algorithm string
	Status    string
	Keyword   string
}

func (f RuleFilter) Match(r *Rule) bool {
	if f.Algorithm != "" && r.Algorithm != f.Algorithm {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Key), k) {
			return false
		}
	}
	return true
}
