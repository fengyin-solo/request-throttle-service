package service

import (
	"sort"
	"time"

	"ratelimiter/internal/model"
	"ratelimiter/pkg/idgen"
)

func (s *Service) CreateAccessRecord(input model.AccessRecord) (*model.AccessRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	a := &model.AccessRecord{
		ID:        idgen.Hex(),
		ClientID:  input.ClientID,
		RuleID:    input.RuleID,
		Path:      input.Path,
		Allowed:   input.Allowed,
		Remaining: input.Remaining,
		LatencyMs: input.LatencyMs,
		CreatedAt: now,
	}
	if err := s.store.CreateAccessRecord(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) GetAccessRecord(id string) (*model.AccessRecord, error) {
	return s.store.GetAccessRecord(id)
}

func (s *Service) ListAccessRecords(filter model.AccessRecordFilter, page, size int) ([]*model.AccessRecord, int, error) {
	all := s.store.ListAccessRecords()
	matched := make([]*model.AccessRecord, 0, len(all))
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
		return []*model.AccessRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
