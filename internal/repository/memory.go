package repository

import (
	"collab-editor/internal/domain"
	"errors"
	"sync"
)

type InMemoryRepo struct {
	store map[string]*domain.Document
	mu    sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		store: make(map[string]*domain.Document),
	}
}

func (r *InMemoryRepo) GetByID(id string) (*domain.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if doc, ok := r.store[id]; ok {
		return doc, nil
	}

	return nil, errors.New("document not found")
}

func (r *InMemoryRepo) Save(doc *domain.Document) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[doc.ID] = doc
	return nil
}
