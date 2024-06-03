package services

import (
	"boilgopher/storage"
	"boilgopher/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler for account opening document upload related
type Book struct {
	storage storage.Storage
	// bookUsecase usecase.IBookUsecase
}

func New(storage storage.Storage) *Book {
	return &Book{
		storage: storage,
	}
}

func (h Book) Book(r *gin.RouterGroup) {
	r.POST(utils.BOOK_URL, h.PostBook) // private
}

func (h Book) PostBook(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "you are not allowed",
	})
}
