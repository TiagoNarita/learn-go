package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tiagobnarita/go_learn/internal/model"
)

var ErrNotFound = errors.New("bookmark not found")

type Repository interface {
	Create(ctx context.Context, b model.Bookmark) (model.Bookmark, error)
	List(ctx context.Context, filter BookmarkFilter) ([]model.Bookmark, int, error)
	GetById(ctx context.Context, id uuid.UUID) (model.Bookmark, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, bookmark model.Bookmark) (model.Bookmark, error)
}

type BookmarkFilter struct {
	Tag    string
	Title  string
	OffSet int
	Limit  int
	Url    string
}
