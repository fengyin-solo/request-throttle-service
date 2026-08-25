package store

import (
	"ratelimiter/internal/model"
)

func (s *MemoryStore) CreateRule(r *model.Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.rules {
		if exist.Key == r.Key {
			return ErrConflict
		}
	}
	s.rules[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRule(id string) (*model.Rule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) GetRuleByKey(key string) (*model.Rule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.rules {
		if r.Key == key {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRules() []*model.Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Rule, 0, len(s.rules))
	for _, r := range s.rules {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRule(r *model.Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[r.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.rules {
		if exist.ID != r.ID && exist.Key == r.Key {
			return ErrConflict
		}
	}
	s.rules[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[id]; !ok {
		return ErrNotFound
	}
	delete(s.rules, id)
	return nil
}
