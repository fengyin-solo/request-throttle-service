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
	existing.Key = input.Key
	existing.Algorithm = input.Algorithm
	existing.Limit = input.Limit
	existing.WindowSec = input.WindowSec
	existing.Status = input.Status
	existing.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateRule(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteRule(id string) error {
	return s.store.DeleteRule(id)
}
