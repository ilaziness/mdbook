package action

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilaziness/mdbook/internal/booklib"
)

func Book(g *gin.Context) {
	bookname := g.Param("bookname")
	if err := booklib.BookExist(bookname); err != nil {
		g.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	g.JSON(200, booklib.Lib[bookname])
}

func BookPage(g *gin.Context) {
	bookname := g.Param("bookname")
	pagepath := g.Param("path")
	if err := booklib.BookExist(bookname); err != nil {
		g.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	g.String(200, booklib.GetPageContent(bookname, pagepath))
}
