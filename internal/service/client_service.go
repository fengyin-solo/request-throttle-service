package service

import (
	"sort"
	"time"

	"ratelimiter/internal/model"
	"ratelimiter/pkg/idgen"
)

func (s *Service) CreateClient(input model.Client) (*model.Client, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	c := &model.Client{
		ID:         idgen.Hex(),
		AppID:      input.AppID,
		Name:       input.Name,
		DailyQuota: input.DailyQuota,
		Status:     input.Status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateClient(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetClient(id string) (*model.Client, error) {
	return s.store.GetClient(id)
}

func (s *Service) GetClientByAppID(appID string) (*model.Client, error) {
	return s.store.GetClientByAppID(appID)
}

func (s *Service) ListClients(filter model.ClientFilter, page, size int) ([]*model.Client, int, error) {
	all := s.store.ListClients()
	matched := make([]*model.Client, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Client{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateClient(id string, input model.Client) (*model.Client, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetClient(id)
	if err != nil {
		return nil, err
	}
	// 在拷贝上应用改动：existing 直接指向 store 中的对象，若先改它再校验，
	// 一旦 AppID 冲突返回 ErrConflict，原对象已被改脏，后续查询会读到半更新内容。
	updated := *existing
	updated.AppID = input.AppID
	updated.Name = input.Name
	updated.DailyQuota = input.DailyQuota
	updated.Status = input.Status
	updated.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateClient(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *Service) DeleteClient(id string) error {
	return s.store.DeleteClient(id)
}

func (s *Service) BanClient(id string) (*model.Client, error) {
	c, err := s.store.GetClient(id)
	if err != nil {
		return nil, err
	}
	if !c.CanTransition(model.ClientStatusSuspended) {
		return nil, model.NewValidationError("status", "当前状态无法封禁")
	}
	c.Status = model.ClientStatusSuspended
	c.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateClient(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) UnbanClient(id string) (*model.Client, error) {
	c, err := s.store.GetClient(id)
	if err != nil {
		return nil, err
	}
	if !c.CanTransition(model.ClientStatusActive) {
		return nil, model.NewValidationError("status", "当前状态无法解封")
	}
	c.Status = model.ClientStatusActive
	c.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateClient(c); err != nil {
		return nil, err
	}
	return c, nil
}
