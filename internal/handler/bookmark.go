package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tiagobnarita/go_learn/internal/handler/dto"
	"github.com/tiagobnarita/go_learn/internal/http/response"
	"github.com/tiagobnarita/go_learn/internal/repository"
	"github.com/tiagobnarita/go_learn/internal/service"
)

type BookmarkHandler struct {
	service service.BookmarkService
}

func (h *BookmarkHandler) Register(router *gin.RouterGroup) {
	router.POST("/", h.Create)
	router.GET("/", h.List)
	router.GET("/:id", h.GetById)
	router.PUT("/:id", h.Update)
	router.DELETE("/:id", h.Delete)
	router.PATCH("/:id", h.Patch)
}

func NewBookmarkHandler(service service.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{
		service: service,
	}
}

func (h *BookmarkHandler) Create(c *gin.Context) {
	var req dto.CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.service.Create(c.Request.Context(), req.ToInput())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, dto.FromDomain(created))
}

func (h *BookmarkHandler) List(c *gin.Context) {
	var req dto.BookmarkPaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	filter := repository.BookmarkFilter{
		Title:  req.Title,
		Tag:    req.Tag,
		Limit:  req.GetLimit(),
		OffSet: req.GetOffset(),
		Url:    req.Url,
	}

	list, total, err := h.service.List(c, filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.NewPageResponse(list, total, req.Page, filter.Limit))
}

func (h *BookmarkHandler) GetById(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	bookmark, err := h.service.GetById(ctx, id)

	if errors.Is(err, repository.ErrNotFound) {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, dto.FromDomain(bookmark))
}

func (h *BookmarkHandler) Delete(ctx *gin.Context) {
	bookMarkId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Delete(ctx, bookMarkId)

	if errors.Is(err, repository.ErrNotFound) {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusAccepted)
}

func (h *BookmarkHandler) Update(ctx *gin.Context) {
	var req dto.CreateBookmarkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	bookmarkId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	update, err := h.service.Update(ctx, bookmarkId, req.ToInput())

	if errors.Is(err, repository.ErrNotFound) {
		response.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, dto.FromDomain(update))
}

func (h *BookmarkHandler) Patch(context *gin.Context) {
	var req dto.BookmarkPatch
	if err := context.ShouldBindJSON(&req); err != nil {
		response.Error(context, http.StatusBadRequest, err.Error())
		return
	}

	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		response.Error(context, http.StatusBadRequest, err.Error())
		return
	}

	update, err := h.service.Patch(context, id, dto.BookmarkPatchToInput(req))

	if errors.Is(err, repository.ErrNotFound) {
		response.Error(context, http.StatusNotFound, err.Error())
		return
	}

	context.JSON(http.StatusOK, dto.FromDomain(update))
}
