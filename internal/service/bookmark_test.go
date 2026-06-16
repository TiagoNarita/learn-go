package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/tiagobnarita/go_learn/internal/model"
	"github.com/tiagobnarita/go_learn/internal/repository"
)

type fakeRepo struct {
	created    model.Bookmark
	listItems  []model.Bookmark
	listTotal  int
	listFilter repository.BookmarkFilter
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

func (f *fakeRepo) List(ctx context.Context, filter repository.BookmarkFilter) ([]model.Bookmark, int, error) {
	f.listFilter = filter
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

func strPtr(s string) *string { return &s }

func TestBookmarkService_List(t *testing.T) {
	t.Run("returns items and total from repo", func(t *testing.T) {
		ctx := context.Background()
		items := []model.Bookmark{{ID: uuid.New()}, {ID: uuid.New()}}
		svc := NewBookmarkService(&fakeRepo{listItems: items, listTotal: 2})

		got, total, err := svc.List(ctx, repository.BookmarkFilter{})
		if err != nil {
			t.Fatalf("List returned error: %v", err)
		}
		if len(got) != 2 || total != 2 {
			t.Fatalf("got %d items total %d, want 2/2", len(got), total)
		}
	})

	t.Run("passes filter through to repo", func(t *testing.T) {
		ctx := context.Background()
		fakeRep := &fakeRepo{}
		svc := NewBookmarkService(fakeRep)
		want := repository.BookmarkFilter{Tag: "go", Title: "gin", Url: "http", Limit: 10, OffSet: 20}

		if _, _, err := svc.List(ctx, want); err != nil {
			t.Fatalf("List returned error: %v", err)
		}
		if fakeRep.listFilter != want {
			t.Errorf("filter passed to repo = %+v, want %+v", fakeRep.listFilter, want)
		}
	})

	t.Run("propagates repo error", func(t *testing.T) {
		ctx := context.Background()
		svc := NewBookmarkService(&fakeRepo{err: errors.New("boom")})

		if _, _, err := svc.List(ctx, repository.BookmarkFilter{}); err == nil {
			t.Fatal("expected error from repo")
		}
	})
}

func TestBookmarkService_Patch(t *testing.T) {
	existing := model.Bookmark{
		ID:    uuid.New(),
		URL:   "old-url",
		Title: "old-title",
		Tags:  []string{"old"},
		Notes: "old-notes",
	}

	tests := []struct {
		name  string
		input PatchBookmarkInput
		want  model.Bookmark
	}{
		{
			name:  "empty patch keeps everything",
			input: PatchBookmarkInput{},
			want:  existing,
		},
		{
			name:  "patches title only",
			input: PatchBookmarkInput{Title: strPtr("new-title")},
			want:  model.Bookmark{ID: existing.ID, URL: "old-url", Title: "new-title", Tags: []string{"old"}, Notes: "old-notes"},
		},
		{
			name:  "patches url only",
			input: PatchBookmarkInput{URL: strPtr("new-url")},
			want:  model.Bookmark{ID: existing.ID, URL: "new-url", Title: "old-title", Tags: []string{"old"}, Notes: "old-notes"},
		},
		{
			name:  "patches tags only",
			input: PatchBookmarkInput{Tags: []string{"a", "b"}},
			want:  model.Bookmark{ID: existing.ID, URL: "old-url", Title: "old-title", Tags: []string{"a", "b"}, Notes: "old-notes"},
		},
		{
			name:  "patches notes only",
			input: PatchBookmarkInput{Notes: strPtr("new-notes")},
			want:  model.Bookmark{ID: existing.ID, URL: "old-url", Title: "old-title", Tags: []string{"old"}, Notes: "new-notes"},
		},
		{
			name:  "patches all fields",
			input: PatchBookmarkInput{URL: strPtr("u"), Title: strPtr("t"), Tags: []string{"x"}, Notes: strPtr("n")},
			want:  model.Bookmark{ID: existing.ID, URL: "u", Title: "t", Tags: []string{"x"}, Notes: "n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			fakeRep := &fakeRepo{getResult: existing}
			svc := NewBookmarkService(fakeRep)

			got, err := svc.Patch(ctx, existing.ID, tt.input)
			if err != nil {
				t.Fatalf("Patch returned error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Patch() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, fakeRep.updated); diff != "" {
				t.Errorf("repo received mismatch (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("when not found", func(t *testing.T) {
		ctx := context.Background()
		svc := NewBookmarkService(&fakeRepo{getErr: repository.ErrNotFound})

		_, err := svc.Patch(ctx, uuid.New(), PatchBookmarkInput{Title: strPtr("x")})
		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("error should be notfound, got %v", err)
		}
	})

	t.Run("propagates update error", func(t *testing.T) {
		ctx := context.Background()
		svc := NewBookmarkService(&fakeRepo{getResult: existing, updateErr: errors.New("boom")})

		if _, err := svc.Patch(ctx, existing.ID, PatchBookmarkInput{Title: strPtr("x")}); err == nil {
			t.Fatal("expected error from repo update")
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
