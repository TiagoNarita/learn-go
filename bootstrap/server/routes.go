package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tiagobnarita/go_learn/internal/handler"
	"github.com/tiagobnarita/go_learn/internal/http/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(middleware.CorsMiddleware())

	basePath := "/bookmarks"

	group := r.Group(basePath)
	//TODO - add swagger

	s.setupBookmarkRouter(group)

	return r
}

func (s *Server) setupBookmarkRouter(group *gin.RouterGroup) {
	bookmarkHandler := handler.NewBookmarkHandler(s.bookmarkService)
	bookmarkHandler.Register(group)
}
