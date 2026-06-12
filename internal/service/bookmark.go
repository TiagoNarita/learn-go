package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tiagobnarita/go_learn/internal/model"
	"github.com/tiagobnarita/go_learn/internal/repository"
)

type CreateBookmarkInput struct {
	URL   string
	Title string
	Tags  []string
	Notes string
}

type BookmarkService interface {
	Create(ctx context.Context, in CreateBookmarkInput) (model.Bookmark, error)
	List(ctx context.Context, limit, offset int) ([]model.Bookmark, int, error)
	GetById(ctx context.Context, id uuid.UUID) (model.Bookmark, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, input CreateBookmarkInput) (model.Bookmark, error)
}

type bookmarkService struct {
	repo repository.Repository
}

func NewBookmarkService(repo repository.Repository) BookmarkService {
	return &bookmarkService{
		repo: repo,
	}
}

func (s *bookmarkService) Create(ctx context.Context, in CreateBookmarkInput) (model.Bookmark, error) {
	bookmark := model.Bookmark{
		ID:        uuid.New(),
		URL:       in.URL,
		Title:     in.Title,
		Tags:      in.Tags,
		Notes:     in.Notes,
		CreatedAt: time.Now().UTC(),
	}
	return s.repo.Create(ctx, bookmark)
}

func (s *bookmarkService) List(ctx context.Context, limit, offset int) ([]model.Bookmark, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *bookmarkService) GetById(ctx context.Context, id uuid.UUID) (model.Bookmark, error) {
	return s.repo.GetById(ctx, id)
}

func (s *bookmarkService) Update(ctx context.Context, id uuid.UUID, input CreateBookmarkInput) (model.Bookmark, error) {
	bookmark, err := s.GetById(ctx, id)
	if err != nil {
		return model.Bookmark{}, err
	}

	updatedBookmark := model.Bookmark{
		ID:    bookmark.ID,
		URL:   input.URL,
		Title: input.Title,
		Tags:  input.Tags,
		Notes: input.Notes,
	}

	return s.repo.Update(ctx, updatedBookmark)
}

func (s *bookmarkService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
