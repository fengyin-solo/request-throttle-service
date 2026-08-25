package store

import (
	"ratelimiter/internal/model"
)

func (s *MemoryStore) CreateBucket(b *model.Bucket) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.buckets {
		if exist.ClientID == b.ClientID && exist.RuleID == b.RuleID {
			return ErrConflict
		}
	}
	s.buckets[b.ID] = b
	return nil
}

func (s *MemoryStore) GetBucket(id string) (*model.Bucket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.buckets[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *b
	return &cp, nil
}

func (s *MemoryStore) GetBucketByClientRule(clientID, ruleID string) (*model.Bucket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, b := range s.buckets {
		if b.ClientID == clientID && b.RuleID == ruleID {
			cp := *b
			return &cp, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListBuckets() []*model.Bucket {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Bucket, 0, len(s.buckets))
	for _, b := range s.buckets {
		cp := *b
		list = append(list, &cp)
	}
	return list
}

func (s *MemoryStore) UpdateBucket(b *model.Bucket) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.buckets[b.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.buckets {
		if exist.ID != b.ID && exist.ClientID == b.ClientID && exist.RuleID == b.RuleID {
			return ErrConflict
		}
	}
	s.buckets[b.ID] = b
	return nil
}

func (s *MemoryStore) DeleteBucket(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.buckets[id]; !ok {
		return ErrNotFound
	}
	delete(s.buckets, id)
	return nil
}
