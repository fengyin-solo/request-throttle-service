package service

import (
	"fmt"
	"time"

	"ratelimiter/internal/model"
	"ratelimiter/pkg/idgen"
)

type CheckResult struct {
	Allowed   bool   `json:"allowed"`
	Remaining int    `json:"remaining"`
	Message   string `json:"message"`
}

func (s *Service) Check(clientID, ruleKey, path string) (*CheckResult, error) {
	client, err := s.store.GetClientByAppID(clientID)
	if err != nil {
		return nil, model.NewValidationError("client_id", "客户端不存在")
	}

	rule, err := s.store.GetRuleByKey(ruleKey)
	if err != nil {
		return nil, model.NewValidationError("rule_key", "规则不存在")
	}

	if client.Status == model.ClientStatusSuspended {
		_ = s.store.CreateAlert(&model.Alert{
			ID:          idgen.Hex(),
			ClientID:    client.ID,
			RuleID:      rule.ID,
			Level:       model.AlertLevelCritical,
			Message:     fmt.Sprintf("客户端 %s 已封禁，请求 %s 被拒绝", client.AppID, path),
			TriggeredAt: time.Now().UTC(),
			CreatedAt:   time.Now().UTC(),
		})
		_ = s.store.CreateAccessRecord(&model.AccessRecord{
			ID:        idgen.Hex(),
			ClientID:  client.ID,
			RuleID:    rule.ID,
			Path:      path,
			Allowed:   false,
			Remaining: 0,
			LatencyMs: 0,
			CreatedAt: time.Now().UTC(),
		})
		return &CheckResult{Allowed: false, Remaining: 0, Message: "客户端已被封禁"}, nil
	}

	now := time.Now().UTC()
	window := time.Duration(rule.WindowSec) * time.Second

	var allowed bool
	var remaining int
	var message string

	switch rule.Algorithm {
	case model.RuleAlgorithmFixedWindow:
		windowStart := now.Truncate(window)
		count := s.countAccessRecords(client.ID, rule.ID, windowStart, now)
		if count >= rule.Limit {
			allowed = false
			remaining = 0
			message = "请求超过固定窗口限流阈值"
		} else {
			allowed = true
			remaining = rule.Limit - count - 1
			message = "ok"
		}
	case model.RuleAlgorithmSlidingWindow:
		windowStart := now.Add(-window)
		count := s.countAccessRecords(client.ID, rule.ID, windowStart, now)
		if count >= rule.Limit {
			allowed = false
			remaining = 0
			message = "请求超过滑动窗口限流阈值"
		} else {
			allowed = true
			remaining = rule.Limit - count - 1
			message = "ok"
		}
	case model.RuleAlgorithmTokenBucket:
		b, err := s.store.GetBucketByClientRule(client.ID, rule.ID)
		if err != nil {
			b = &model.Bucket{
				ID:         idgen.Hex(),
				ClientID:   client.ID,
				RuleID:     rule.ID,
				Tokens:     float64(rule.Limit),
				LastRefill: now,
				UpdatedAt:  now,
			}
			if err := s.store.CreateBucket(b); err != nil {
				return nil, err
			}
		}
		elapsed := now.Sub(b.LastRefill).Seconds()
		refillRate := float64(rule.Limit) / float64(rule.WindowSec)
		b.Tokens += elapsed * refillRate
		if b.Tokens > float64(rule.Limit) {
			b.Tokens = float64(rule.Limit)
		}
		b.LastRefill = now
		if b.Tokens >= 1 {
			b.Tokens--
			allowed = true
			remaining = int(b.Tokens)
			message = "ok"
		} else {
			allowed = false
			remaining = 0
			message = "令牌不足"
		}
		b.UpdatedAt = now
		if err := s.store.UpdateBucket(b); err != nil {
			return nil, err
		}
	}

	latencyMs := 1
	_ = s.store.CreateAccessRecord(&model.AccessRecord{
		ID:        idgen.Hex(),
		ClientID:  client.ID,
		RuleID:    rule.ID,
		Path:      path,
		Allowed:   allowed,
		Remaining: remaining,
		LatencyMs: latencyMs,
		CreatedAt: now,
	})

	return &CheckResult{Allowed: allowed, Remaining: remaining, Message: message}, nil
}

func (s *Service) countAccessRecords(clientID, ruleID string, start, end time.Time) int {
	all := s.store.ListAccessRecords()
	count := 0
	for _, a := range all {
		if a.ClientID == clientID && a.RuleID == ruleID {
			if !a.CreatedAt.Before(start) && !a.CreatedAt.After(end) {
				count++
			}
		}
	}
	return count
}
