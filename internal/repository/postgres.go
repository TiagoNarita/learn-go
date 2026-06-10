package repository

import (
	"context"
	"errors"
	"fmt"

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

func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]model.Bookmark, int, error) {
	total, err := r.countBookmarkTotal(ctx)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return make([]model.Bookmark, 0), 0, nil
	}

	const q = `
		SELECT id, url, title, tags, notes, created_at
		FROM bookmarks
		ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, q, limit, offset)
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

func (r *PostgresRepository) countBookmarkTotal(ctx context.Context) (int, error) {
	var total int
	const countQuery = `SELECT COUNT(*) FROM bookmarks`
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
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
