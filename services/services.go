package services

import (
	"fmt"
	"net/url"

	"github.com/HanjuraNivan/url-shortener/api"
	urlshortenersqlc "github.com/HanjuraNivan/url-shortener/sqlc/urlshortenerqueries"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/thanhpk/randstr"
)

type Services struct {
	sqlcQueries *urlshortenersqlc.Queries
}

func generateErrorMessage(message string) map[string]string {
	return map[string]string{
		"err": message,
	}
}

func (s *Services) RedirectToOriginalUrl(c *gin.Context, shortenedURL string) {
	if shortenedURL == "" {
		c.JSON(500, generateErrorMessage("Short Url needs to be provided to redirect!"))
		return
	}
	fullUrl, err := s.sqlcQueries.GetFullURLFromShortURL(c, shortenedURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(404, generateErrorMessage("No full url found for this short URL!"))
		} else {
			c.JSON(500, generateErrorMessage(err.Error()))
		}
		return
	}
	if fullUrl == "" {
		c.JSON(404, generateErrorMessage("Full URL is empty for this short URL"))
		return
	}
	fmt.Println("Redirecting you to the original URL ....")
	err = s.sqlcQueries.UpdateClickedTimes(c, shortenedURL)
	if err != nil {
		c.JSON(429, generateErrorMessage("Could not update the count for clicked times"))
		return
	}
	c.Redirect(302, fullUrl)
}

func (s *Services) GetOriginalURL(c *gin.Context, params api.GetOriginalURLParams) {
	if params.ShortUrl == "" {
		c.JSON(500, generateErrorMessage("Short URL is needed for getting full URL"))
		return
	}
	fullUrl, err := s.sqlcQueries.GetFullURLFromShortURL(c, params.ShortUrl)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(404, generateErrorMessage("No full URL found for this short URL"))
		} else {
			c.JSON(500, err)
		}
		return
	}
	if fullUrl == "" {
		c.JSON(404, generateErrorMessage("Full URL is empty for this short URL"))
		return
	}
	c.JSON(200, struct {
		FullURL string
	}{
		FullURL: fullUrl,
	})
}

func (s *Services) ShortenURL(c *gin.Context) {
	var bodyReq api.ShortenURLJSONRequestBody
	err := c.ShouldBindJSON(&bodyReq)
	if err != nil {
		c.JSON(500, generateErrorMessage("Failed to bind the body!"))
		return
	}
	if bodyReq.Url == "" {
		c.JSON(404, generateErrorMessage("Please provide a full URL to shorten"))
		return
	}
	_, err = url.ParseRequestURI(bodyReq.Url)
	if err != nil {
		c.JSON(500, err)
		return
	}
	_, err = s.sqlcQueries.GetFullURL(c, bodyReq.Url)
	if err != nil && err != pgx.ErrNoRows {
		c.JSON(500, err)
		return
	}
	createdShortCode, err := validateAndCreateShortCode(c, s.sqlcQueries)
	if err != nil {
		c.JSON(500, err)
		return
	}
	if createdShortCode == "" {
		c.JSON(404, generateErrorMessage("Could not create short code after 5 tries! please try again"))
		return
	}
	_, err = s.sqlcQueries.CreateURLMapping(c, urlshortenersqlc.CreateURLMappingParams{
		Fullurl:  bodyReq.Url,
		Shorturl: createdShortCode,
	})
	if err != nil {
		c.JSON(422, generateErrorMessage("Could not insert the url mapping!"))
		return
	}
	c.JSON(201, struct {
		ShortURL string
	}{
		ShortURL: createdShortCode,
	})
}

func (s *Services) GetUrlClicks(c *gin.Context, shortenedURL string) {
	if shortenedURL == "" {
		c.JSON(409, generateErrorMessage("No short url found!"))
		return
	}
	count, err := s.sqlcQueries.GetClickedTimes(c, shortenedURL)
	if err != nil {
		c.JSON(500, generateErrorMessage(err.Error()))
		return
	}
	c.JSON(200, struct {
		Count int
	}{
		Count: (int(count.Int32)),
	})
}

func validateAndCreateShortCode(c *gin.Context, sqlcQueries *urlshortenersqlc.Queries) (string, error) {
	retries := 0
	for retries < 5 {
		shortCode := randstr.Base62(7)
		_, err := sqlcQueries.GetShortURL(c, shortCode)
		if err != nil {
			if err == pgx.ErrNoRows {
				return shortCode, nil
			} else {
				return "", err
			}
		}
		retries++
	}
	return "", nil
}

func NewServices(sqlcQueries *urlshortenersqlc.Queries) *Services {
	return &Services{
		sqlcQueries: sqlcQueries,
	}
}
