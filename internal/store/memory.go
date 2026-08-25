package store

import (
	"sync"

	"ratelimiter/internal/model"
)

type MemoryStore struct {
	mu            sync.RWMutex
	rules         map[string]*model.Rule
	clients       map[string]*model.Client
	accessRecords map[string]*model.AccessRecord
	buckets       map[string]*model.Bucket
	alerts        map[string]*model.Alert
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rules:         make(map[string]*model.Rule),
		clients:       make(map[string]*model.Client),
		accessRecords: make(map[string]*model.AccessRecord),
		buckets:       make(map[string]*model.Bucket),
		alerts:        make(map[string]*model.Alert),
	}
}

var _ Store = (*MemoryStore)(nil)
