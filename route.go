package mdbook

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilaziness/mdbook/internal/action"
)

func initRoute() {
	BookServer.Engine.GET("/ping", func(g *gin.Context) {
		g.String(http.StatusOK, "pong")
	})
	BookServer.Engine.GET("/book/:bookname", action.Book)
	BookServer.Engine.GET("/book/:bookname/*path", action.BookPage)
}
