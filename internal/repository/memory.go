package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/tiagobnarita/go_learn/internal/model"
)

type InMemoryRepository struct {
	mu    sync.RWMutex
	items map[uuid.UUID]model.Bookmark
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items: make(map[uuid.UUID]model.Bookmark),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, b model.Bookmark) (model.Bookmark, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[b.ID] = b
	return r.items[b.ID], nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]model.Bookmark, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.Bookmark, 0, len(r.items))
	for _, b := range r.items {
		out = append(out, b)
	}
	return out, nil
}
