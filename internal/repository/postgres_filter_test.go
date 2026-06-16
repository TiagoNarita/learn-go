package repository

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateQueryForFilter(t *testing.T) {
	tests := []struct {
		name      string
		filter    BookmarkFilter
		wantConds []string
		wantArgs  []any
		wantNext  int
	}{
		{
			name:      "empty filter has no conditions",
			filter:    BookmarkFilter{},
			wantConds: []string{},
			wantArgs:  []any{},
			wantNext:  1,
		},
		{
			name:      "tag uses equality",
			filter:    BookmarkFilter{Tag: "go"},
			wantConds: []string{"tags = $1"},
			wantArgs:  []any{"go"},
			wantNext:  2,
		},
		{
			name:      "title uses ILIKE with wildcards",
			filter:    BookmarkFilter{Title: "gin"},
			wantConds: []string{"title ILIKE $1"},
			wantArgs:  []any{"%gin%"},
			wantNext:  2,
		},
		{
			name:      "url uses ILIKE with wildcards",
			filter:    BookmarkFilter{Url: "http"},
			wantConds: []string{"url ILIKE $1"},
			wantArgs:  []any{"%http%"},
			wantNext:  2,
		},
		{
			name:      "all filters increment placeholders in order",
			filter:    BookmarkFilter{Tag: "go", Title: "gin", Url: "http"},
			wantConds: []string{"tags = $1", "title ILIKE $2", "url ILIKE $3"},
			wantArgs:  []any{"go", "%gin%", "%http%"},
			wantNext:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conds, args, next := createQueryFoFilter(tt.filter)

			if diff := cmp.Diff(tt.wantConds, conds); diff != "" {
				t.Errorf("conds mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantArgs, args); diff != "" {
				t.Errorf("args mismatch (-want +got):\n%s", diff)
			}
			if next != tt.wantNext {
				t.Errorf("next placeholder = %d, want %d", next, tt.wantNext)
			}
		})
	}
}

func TestApplyFilter(t *testing.T) {
	base := "SELECT id FROM bookmarks"

	t.Run("no conditions skips WHERE and appends paging", func(t *testing.T) {
		query, args := applyFilter(nil, base, 1, []any{}, BookmarkFilter{Limit: 10, OffSet: 20})

		if strings.Contains(query, "WHERE") {
			t.Errorf("query should not contain WHERE: %q", query)
		}
		if !strings.Contains(query, "ORDER BY created_at DESC LIMIT $1 OFFSET $2") {
			t.Errorf("query missing paging clause: %q", query)
		}
		if diff := cmp.Diff([]any{10, 20}, args); diff != "" {
			t.Errorf("args mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("conditions joined with AND and paging uses next placeholders", func(t *testing.T) {
		conds := []string{"tags = $1", "title ILIKE $2"}
		query, args := applyFilter(conds, base, 3, []any{"go", "%gin%"}, BookmarkFilter{Limit: 5, OffSet: 0})

		if !strings.Contains(query, "WHERE tags = $1 AND title ILIKE $2") {
			t.Errorf("query missing joined conditions: %q", query)
		}
		if !strings.Contains(query, "LIMIT $3 OFFSET $4") {
			t.Errorf("paging should continue from next placeholder: %q", query)
		}
		if diff := cmp.Diff([]any{"go", "%gin%", 5, 0}, args); diff != "" {
			t.Errorf("args mismatch (-want +got):\n%s", diff)
		}
	})
}