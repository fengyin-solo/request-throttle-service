// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"ratelimiter/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Rule
	CreateRule(r *model.Rule) error
	GetRule(id string) (*model.Rule, error)
	GetRuleByKey(key string) (*model.Rule, error)
	ListRules() []*model.Rule
	UpdateRule(r *model.Rule) error
	DeleteRule(id string) error

	// Client
	CreateClient(c *model.Client) error
	GetClient(id string) (*model.Client, error)
	GetClientByAppID(appID string) (*model.Client, error)
	ListClients() []*model.Client
	UpdateClient(c *model.Client) error
	DeleteClient(id string) error

	// AccessRecord（只读追加）
	CreateAccessRecord(a *model.AccessRecord) error
	GetAccessRecord(id string) (*model.AccessRecord, error)
	ListAccessRecords() []*model.AccessRecord

	// Bucket
	CreateBucket(b *model.Bucket) error
	GetBucket(id string) (*model.Bucket, error)
	GetBucketByClientRule(clientID, ruleID string) (*model.Bucket, error)
	ListBuckets() []*model.Bucket
	UpdateBucket(b *model.Bucket) error
	DeleteBucket(id string) error

	// Alert
	CreateAlert(a *model.Alert) error
	GetAlert(id string) (*model.Alert, error)
	ListAlerts() []*model.Alert
	UpdateAlert(a *model.Alert) error
	DeleteAlert(id string) error
}
