package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/tiagobnarita/go_learn/internal/model"
)

func TestInMemoryRepository_Create(t *testing.T) {
	repo := NewInMemoryRepository()
	ctx := context.Background()
	want := createBookmark()

	got, err := repo.Create(ctx, want)

	if err != nil {
		t.Fatalf("Failed to create bookmark: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Create() mismatch (-want +got):\n%s", diff)
	}
}

func createBookmark() model.Bookmark {
	return model.Bookmark{
		ID:        uuid.New(),
		URL:       "url",
		Title:     "title",
		Tags:      []string{"tag1", "tag2"},
		Notes:     "notes",
		CreatedAt: time.Now().UTC(),
	}
}

func TestInMemoryRepository_List(t *testing.T) {
	t.Run("empty repo returns empty slice", func(t *testing.T) {
		repo := NewInMemoryRepository()
		ctx := context.Background()

		got, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("List returned error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("got %d items, want 0", len(got))
		}
	})

	t.Run("returns all created bookmarks", func(t *testing.T) {
		repo := NewInMemoryRepository()
		ctx := context.Background()
		b1 := createBookmark()
		b2 := createBookmark()
		if _, err := repo.Create(ctx, b1); err != nil {
			t.Fatalf("setup Create b1: %v", err)
		}
		if _, err := repo.Create(ctx, b2); err != nil {
			t.Fatalf("setup Create b2: %v", err)
		}

		got, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("List returned error: %v", err)
		}

		want := []model.Bookmark{b1, b2}
		less := func(a, b model.Bookmark) bool { return a.ID.String() < b.ID.String() }
		if diff := cmp.Diff(want, got, cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("List() mismatch (-want +got):\n%s", diff)
		}
	})
}
