package store

import (
	"ratelimiter/internal/model"
)

func (s *MemoryStore) CreateClient(c *model.Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.clients {
		if exist.AppID == c.AppID {
			return ErrConflict
		}
	}
	s.clients[c.ID] = c
	return nil
}

func (s *MemoryStore) GetClient(id string) (*model.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.clients[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *c
	return &cp, nil
}

func (s *MemoryStore) GetClientByAppID(appID string) (*model.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.clients {
		if c.AppID == appID {
			cp := *c
			return &cp, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListClients() []*model.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Client, 0, len(s.clients))
	for _, c := range s.clients {
		cp := *c
		list = append(list, &cp)
	}
	return list
}

func (s *MemoryStore) UpdateClient(c *model.Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clients[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.clients {
		if exist.ID != c.ID && exist.AppID == c.AppID {
			return ErrConflict
		}
	}
	s.clients[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteClient(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clients[id]; !ok {
		return ErrNotFound
	}
	delete(s.clients, id)
	return nil
}
