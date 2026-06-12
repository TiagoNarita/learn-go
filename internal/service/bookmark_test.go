package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/tiagobnarita/go_learn/internal/model"
	"github.com/tiagobnarita/go_learn/internal/repository"
)

type fakeRepo struct {
	created    model.Bookmark
	listItems  []model.Bookmark
	listTotal  int
	err        error
	getResult  model.Bookmark
	getErr     error
	deletedErr error
	updated    model.Bookmark
	updateErr  error
}

func (f *fakeRepo) Create(ctx context.Context, b model.Bookmark) (model.Bookmark, error) {
	f.created = b
	return b, f.err
}

func (f *fakeRepo) List(ctx context.Context, limit, offset int) ([]model.Bookmark, int, error) {
	return f.listItems, f.listTotal, f.err
}

func (f *fakeRepo) GetById(ctx context.Context, id uuid.UUID) (model.Bookmark, error) {
	return f.getResult, f.getErr
}

func (f *fakeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return f.deletedErr
}

func (f *fakeRepo) Update(ctx context.Context, bookmark model.Bookmark) (model.Bookmark, error) {
	f.updated = bookmark
	return bookmark, f.updateErr
}
func createBookmarkInput() CreateBookmarkInput {
	return CreateBookmarkInput{
		URL:   "url",
		Title: "title",
		Tags:  []string{"tag1", "tag2"},
		Notes: "notes",
	}
}

func TestBookmarkService_Create(t *testing.T) {
	ctx := context.Background()
	fakeRep := &fakeRepo{}

	svc := NewBookmarkService(fakeRep)
	created, err := svc.Create(ctx, createBookmarkInput())
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if created.ID == uuid.Nil {
		t.Fatalf("id should not be nil")
	}

	if created.URL != "url" {
		t.Fatalf("Create returned url: %v", created.URL)
	}

	if fakeRep.created.URL != "url" {
		t.Errorf("repo got URL: %v", fakeRep.created.URL)
	}

	if fakeRep.created.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should not be zero")
	}
}

func TestBookmarkService_Create_Err(t *testing.T) {
	ctx := context.Background()
	fakeRep := &fakeRepo{err: errors.New("err")}
	service := NewBookmarkService(fakeRep)

	_, err := service.Create(ctx, createBookmarkInput())

	if err == nil {
		t.Fatalf("expected error from erro")
	}
}

func TestBookmarkService_GetById(t *testing.T) {
	t.Run("When found", func(t *testing.T) {
		ctx := context.Background()
		id := uuid.New()
		want := model.Bookmark{ID: id, URL: "url", Title: "title", Tags: []string{"tag1"}}
		service := NewBookmarkService(&fakeRepo{getResult: want})

		byId, err := service.GetById(ctx, id)
		if err != nil {
			t.Fatalf("should not receive error")
		}

		if byId.ID != id {
			t.Fatalf("id should not be different")
		}
	})

	t.Run("when not found", func(t *testing.T) {
		ctx := context.Background()
		service := NewBookmarkService(&fakeRepo{getErr: repository.ErrNotFound})
		_, err := service.GetById(ctx, uuid.New())

		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("error should be notfound")
		}
	})
}

func TestBookmarkService_Update(t *testing.T) {
	t.Run("when updated", func(t *testing.T) {
		ctx := context.Background()
		id := uuid.New()
		existing := model.Bookmark{ID: id, URL: "old", Title: "old", Tags: []string{"old"}}
		fakeRep := &fakeRepo{getResult: existing}
		service := NewBookmarkService(fakeRep)

		update, err := service.Update(ctx, id, createBookmarkInput())
		if err != nil {
			t.Fatalf("should not receive error: %v", err)
		}

		if update.ID != id {
			t.Errorf("id should be kept: got %v want %v", update.ID, id)
		}

		if update.URL != "url" {
			t.Errorf("url should be updated: got %v", update.URL)
		}

		if fakeRep.updated.URL != "url" {
			t.Errorf("repo got URL: %v", fakeRep.updated.URL)
		}
	})

	t.Run("when not found", func(t *testing.T) {
		ctx := context.Background()
		service := NewBookmarkService(&fakeRepo{getErr: repository.ErrNotFound})
		_, err := service.Update(ctx, uuid.New(), createBookmarkInput())

		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("error should be notfound, got %v", err)
		}
	})
}

func TestBookmarkService_Delete(t *testing.T) {
	t.Run("when deleted", func(t *testing.T) {
		ctx := context.Background()
		service := NewBookmarkService(&fakeRepo{})

		err := service.Delete(ctx, uuid.New())
		if err != nil {
			t.Fatalf("should not receive error: %v", err)
		}
	})

	t.Run("when not found", func(t *testing.T) {
		ctx := context.Background()
		service := NewBookmarkService(&fakeRepo{deletedErr: repository.ErrNotFound})

		err := service.Delete(ctx, uuid.New())
		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("error should be notfound, got %v", err)
		}
	})
}
