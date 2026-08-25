package service

import (
	"sort"
	"time"

	"ratelimiter/internal/model"
	"ratelimiter/pkg/idgen"
)

func (s *Service) CreateAlert(input model.Alert) (*model.Alert, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	a := &model.Alert{
		ID:          idgen.Hex(),
		ClientID:    input.ClientID,
		RuleID:      input.RuleID,
		Level:       input.Level,
		Message:     input.Message,
		TriggeredAt: now,
		CreatedAt:   now,
	}
	if err := s.store.CreateAlert(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) GetAlert(id string) (*model.Alert, error) {
	return s.store.GetAlert(id)
}

func (s *Service) ListAlerts(filter model.AlertFilter, page, size int) ([]*model.Alert, int, error) {
	all := s.store.ListAlerts()
	matched := make([]*model.Alert, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Alert{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateAlert(id string, input model.Alert) (*model.Alert, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetAlert(id)
	if err != nil {
		return nil, err
	}
	existing.ClientID = input.ClientID
	existing.RuleID = input.RuleID
	existing.Level = input.Level
	existing.Message = input.Message
	existing.TriggeredAt = input.TriggeredAt
	if err := s.store.UpdateAlert(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteAlert(id string) error {
	return s.store.DeleteAlert(id)
}
