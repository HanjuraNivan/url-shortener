package services

import (
	"errors"

	urlshortenersqlc "github.com/HanjuraNivan/url-shortener/sqlc/urlshortenerqueries"
	"github.com/gin-gonic/gin"
)

type Services struct {
	sqlcQueries *urlshortenersqlc.Queries
}
type UrlShorteningRequest struct {
	FullUrl string
}

func (s *Services) GetOriginalURL(c *gin.Context) {
	var urlRequest UrlShorteningRequest
	err := c.ShouldBindJSON(&urlRequest)
	if err != nil {
		c.JSON(500, errors.New("failed to parse request!"))
		return
	}

}
