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
	List(ctx context.Context, limit, offset int) ([]model.Bookmark, int, error)
	GetById(ctx context.Context, id uuid.UUID) (model.Bookmark, error)
}
