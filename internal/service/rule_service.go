package service

import (
	"sort"
	"time"

	"ratelimiter/internal/model"
	"ratelimiter/pkg/idgen"
)

func (s *Service) CreateRule(input model.Rule) (*model.Rule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	r := &model.Rule{
		ID:        idgen.Hex(),
		Key:       input.Key,
		Algorithm: input.Algorithm,
		Limit:     input.Limit,
		WindowSec: input.WindowSec,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetRule(id string) (*model.Rule, error) {
	return s.store.GetRule(id)
}

func (s *Service) ListRules(filter model.RuleFilter, page, size int) ([]*model.Rule, int, error) {
	all := s.store.ListRules()
	matched := make([]*model.Rule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Rule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRule(id string, input model.Rule) (*model.Rule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetRule(id)
	if err != nil {
		return nil, err
	}
	// 构造副本再交给 store 校验，避免在冲突检查未通过前就原地改动存储中的原规则。
	updated := &model.Rule{
		ID:        existing.ID,
		Key:       input.Key,
		Algorithm: input.Algorithm,
		Limit:     input.Limit,
		WindowSec: input.WindowSec,
		Status:    input.Status,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.store.UpdateRule(updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeleteRule(id string) error {
	return s.store.DeleteRule(id)
}
