package service

import (
	"math"
	"sort"
	"time"

	"ratelimiter/internal/model"
	"ratelimiter/pkg/idgen"
)

func (s *Service) CreateBucket(input model.Bucket) (*model.Bucket, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	b := &model.Bucket{
		ID:         idgen.Hex(),
		ClientID:   input.ClientID,
		RuleID:     input.RuleID,
		Tokens:     input.Tokens,
		LastRefill: now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateBucket(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) GetBucket(id string) (*model.Bucket, error) {
	return s.store.GetBucket(id)
}

func (s *Service) GetBucketByClientRule(clientID, ruleID string) (*model.Bucket, error) {
	return s.store.GetBucketByClientRule(clientID, ruleID)
}

func (s *Service) ListBuckets(filter model.BucketFilter, page, size int) ([]*model.Bucket, int, error) {
	all := s.store.ListBuckets()
	matched := make([]*model.Bucket, 0, len(all))
	for _, b := range all {
		if filter.Match(b) {
			matched = append(matched, b)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].UpdatedAt.After(matched[j].UpdatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Bucket{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateBucket(id string, input model.Bucket) (*model.Bucket, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetBucket(id)
	if err != nil {
		return nil, err
	}
	// 在副本上组装新值，避免在 store 冲突校验通过前改动共享的已有记录。
	// 否则冲突返回后，原桶的 RuleID/Tokens 等字段已被改掉。
	updated := *existing
	updated.ClientID = input.ClientID
	updated.RuleID = input.RuleID
	updated.Tokens = input.Tokens
	updated.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateBucket(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *Service) DeleteBucket(id string) error {
	return s.store.DeleteBucket(id)
}

func (s *Service) RefillBucket(id string, tokens float64) (*model.Bucket, error) {
	if tokens < 0 {
		return nil, model.NewValidationError("tokens", "补充令牌数不能为负数")
	}
	b, err := s.store.GetBucket(id)
	if err != nil {
		return nil, err
	}
	b.Tokens = math.Min(b.Tokens+tokens, float64(1<<53))
	b.LastRefill = time.Now().UTC()
	b.UpdatedAt = b.LastRefill
	if err := s.store.UpdateBucket(b); err != nil {
		return nil, err
	}
	return b, nil
}
