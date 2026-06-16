package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/tiagobnarita/go_learn/internal/model"
	"github.com/tiagobnarita/go_learn/internal/service"
)

type CreateBookmarkRequest struct {
	URL   string   `json:"url" binding:"required,url"`
	Title string   `json:"title" binding:"required,min=1,max=200"`
	Tags  []string `json:"tags" binding:"required,dive,min=1,max=30"`
	Notes string   `json:"notes" binding:"max=200"`
}

func (r CreateBookmarkRequest) ToInput() service.CreateBookmarkInput {
	return service.CreateBookmarkInput{
		URL:   r.URL,
		Title: r.Title,
		Tags:  r.Tags,
		Notes: r.Notes,
	}
}

type BookmarkResponse struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
}

func FromDomain(b model.Bookmark) BookmarkResponse {
	return BookmarkResponse{
		ID:        b.ID,
		URL:       b.URL,
		Title:     b.Title,
		Tags:      b.Tags,
		Notes:     b.Notes,
		CreatedAt: b.CreatedAt,
	}
}

type BookmarkPaginationRequest struct {
	Page  int    `form:"page" json:"page" binding:"min=1"`
	Limit int    `form:"limit" json:"limit" binding:"min=1,max=100"`
	Sort  string `form:"sort" json:"sort"` // e.g., "created_at_desc"
	Tag   string `form:"tag" json:"tag"`
	Title string `form:"title" json:"title"`
	Url   string `form:"url" json:"url"`
}

func (p *BookmarkPaginationRequest) GetOffset() int {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

func (p *BookmarkPaginationRequest) GetLimit() int {
	if p.Limit <= 0 {
		return 10
	}
	return p.Limit
}

type PageResponse[T any] struct {
	Items       []T `json:"items"`
	TotalItems  int `json:"total_items"`
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	Limit       int `json:"limit"`
}

func NewPageResponse[T any](items []T, totalItems int, page, limit int) PageResponse[T] {
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit != 0 {
		totalPages++
	}

	if totalPages == 0 {
		totalPages = 1
	}

	return PageResponse[T]{
		Items:       items,
		TotalItems:  totalItems,
		CurrentPage: page,
		TotalPages:  totalPages,
		Limit:       limit,
	}
}

type BookmarkPatch struct {
	URL   *string  `json:"url" binding:"omitempty,url"`
	Title *string  `json:"title" binding:"omitempty,min=1,max=200"`
	Tags  []string `json:"tags" binding:"omitempty,dive,min=1,max=30"`
	Notes *string  `json:"notes" binding:"omitempty,max=200"`
}

func BookmarkPatchToInput(patch BookmarkPatch) service.PatchBookmarkInput {
	return service.PatchBookmarkInput{
		URL:   patch.URL,
		Title: patch.Title,
		Tags:  patch.Tags,
		Notes: patch.Notes,
	}
}
