package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tiagobnarita/go_learn/internal/model"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, b model.Bookmark) (model.Bookmark, error) {
	const q = `
		INSERT INTO bookmarks (id, url, title, tags, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, url, title, tags, notes, created_at
	`
	var out model.Bookmark
	err := r.pool.QueryRow(ctx, q, b.ID, b.URL, b.Title, b.Tags, b.Notes, b.CreatedAt).Scan(
		&out.ID, &out.URL, &out.Title, &out.Tags, &out.Notes, &out.CreatedAt,
	)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("postgres create bookmark: %w", err)
	}
	return out, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter BookmarkFilter) ([]model.Bookmark, int, error) {
	total, err := r.countBookmarkTotal(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return make([]model.Bookmark, 0), 0, nil
	}

	conds, args, i := createQueryFoFilter(filter)

	query := `
		SELECT id, url, title, tags, notes, created_at
		FROM bookmarks
	`
	query, args = applyFilter(conds, query, i, args, filter)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("postgres list bookmarks: %w", err)
	}
	defer rows.Close()

	out := make([]model.Bookmark, 0)
	for rows.Next() {
		var b model.Bookmark
		if err := rows.Scan(&b.ID, &b.URL, &b.Title, &b.Tags, &b.Notes, &b.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("postgres scan bookmark: %w", err)
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("postgres rows iteration: %w", err)
	}
	return out, total, nil
}

func applyFilter(conds []string, query string, i int, args []any, filter BookmarkFilter) (string, []any) {
	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d "+
		"OFFSET $%d", i, i+1)

	args = append(args, filter.Limit)
	args = append(args, filter.OffSet)
	return query, args
}

func createQueryFoFilter(filter BookmarkFilter) ([]string, []any, int) {
	conds := []string{}
	args := []any{}
	i := 1

	if filter.Tag != "" {
		conds = append(conds, fmt.Sprintf("tags = $%d", i))
		args = append(args, filter.Tag)
		i++
	}

	if filter.Title != "" {
		conds = append(conds, fmt.Sprintf("title ILIKE $%d", i))
		args = append(args, "%"+filter.Title+"%")
		i++
	}

	if filter.Url != "" {
		conds = append(conds, fmt.Sprintf("url ILIKE $%d", i))
		args = append(args, "%"+filter.Url+"%")
		i++
	}
	return conds, args, i
}

func (r *PostgresRepository) countBookmarkTotal(ctx context.Context, filter BookmarkFilter) (int, error) {
	var total int
	conds, args, _ := createQueryFoFilter(filter)
	query := `SELECT COUNT(*) FROM bookmarks`
	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}

	if err := r.pool.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("postgres count bookmarks: %w", err)
	}
	return total, nil
}

func (r *PostgresRepository) GetById(ctx context.Context, id uuid.UUID) (model.Bookmark, error) {
	const q = `
		SELECT id, url, title, tags, notes, created_at
		FROM bookmarks
		WHERE bookmarks.id = $1
	`
	var out model.Bookmark
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&out.ID, &out.URL, &out.Title, &out.Tags, &out.Notes, &out.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.Bookmark{}, ErrNotFound
	}

	if err != nil {
		return model.Bookmark{}, fmt.Errorf("postgress get error: %w", err)
	}

	return out, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `
		DELETE FROM bookmarks 
		WHERE bookmarks.id = $1
	`

	exec, err := r.pool.Exec(ctx, q, id)

	if err != nil {
		return err
	}

	if exec.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, bookmark model.Bookmark) (model.Bookmark, error) {
	const q = `
		UPDATE bookmarks Set url = $1, title = $2, tags = $3, notes = $4 
		WHERE id = $5
		RETURNING id, url, title, tags, notes, created_at
	`
	var out model.Bookmark
	err := r.pool.QueryRow(ctx, q, bookmark.URL, bookmark.Title, bookmark.Tags, bookmark.Notes, bookmark.ID).Scan(
		&out.ID, &out.URL, &out.Title, &out.Tags, &out.Notes, &out.CreatedAt,
	)

	if err != nil {
		return model.Bookmark{}, fmt.Errorf("postgres create bookmark: %w", err)
	}
	return out, nil
}
