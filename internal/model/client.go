package model

import (
	"strings"
	"time"
)

const (
	ClientStatusActive    = "active"
	ClientStatusSuspended = "suspended"
)

type Client struct {
	ID         string    `json:"id"`
	AppID      string    `json:"app_id"`
	Name       string    `json:"name"`
	DailyQuota int64     `json:"daily_quota"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *Client) Validate() error {
	c.AppID = strings.TrimSpace(c.AppID)
	c.Name = strings.TrimSpace(c.Name)
	if c.AppID == "" {
		return NewValidationError("app_id", "AppID 不能为空")
	}
	if c.Name == "" {
		return NewValidationError("name", "客户端名称不能为空")
	}
	if c.Status == "" {
		c.Status = ClientStatusActive
	}
	if c.Status != ClientStatusActive && c.Status != ClientStatusSuspended {
		return NewValidationError("status", "客户端状态不合法")
	}
	return nil
}

func (c *Client) CanTransition(to string) bool {
	if c.Status == ClientStatusActive && to == ClientStatusSuspended {
		return true
	}
	if c.Status == ClientStatusSuspended && to == ClientStatusActive {
		return true
	}
	return false
}

type ClientFilter struct {
	Status  string
	Keyword string
}

func (f ClientFilter) Match(c *Client) bool {
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.AppID), k) && !strings.Contains(strings.ToLower(c.Name), k) {
			return false
		}
	}
	return true
}
