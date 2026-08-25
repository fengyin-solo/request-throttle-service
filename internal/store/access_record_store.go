package store

import (
	"ratelimiter/internal/model"
)

func (s *MemoryStore) CreateAccessRecord(a *model.AccessRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accessRecords[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAccessRecord(id string) (*model.AccessRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.accessRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *a
	return &cp, nil
}

func (s *MemoryStore) ListAccessRecords() []*model.AccessRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AccessRecord, 0, len(s.accessRecords))
	for _, a := range s.accessRecords {
		cp := *a
		list = append(list, &cp)
	}
	return list
}
