package api_service

import (
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"golangAPI/db_service"
	"golangAPI/pojo"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

// POST LongURL and Expired data
func PostUrl(c *gin.Context) {
	url := pojo.URL{}
	c.BindJSON(&url)

	data := db_service.QueryUrl(url.Long_URL)
	if len(data.Long_URL) == 0 {
		url.Id = GenerateShortUrl()
		url.Short_URL = "http://localhost:8000/" + url.Id
		db_service.InsertURL(url.Long_URL, url.Id, url.Short_URL, url.ExpiredDate)
	} else {
		url = data
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       url.Id,
		"shortUrl": url.Short_URL,
	})
}

// GET shortUrl and redirect to original URL
func GetUrl(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request

	url_id := c.Param("url_id")
	data := db_service.QueryId(url_id)
	expiration := ExpireData(data.ExpiredDate)
	if len(data.Id) != 0 && expiration {
		c.JSON(http.StatusNotAcceptable, "The URL is expired.")
	} else if len(data.Id) != 0 && !expiration {
		http.Redirect(w, r, data.Long_URL, 302)
	} else {
		c.JSON(http.StatusNotAcceptable, "The short URL is not exist.")
	}
}

// Create shorter URL
func GenerateShortUrl() string {
	var return_id string = ""
	const base = 36
	size := big.NewInt(base)
	n := make([]byte, 6)
	for i, _ := range n {
		c, _ := rand.Int(rand.Reader, size)
		n[i] = strconv.FormatInt(c.Int64(), base)[0]
	}
	data := db_service.QueryId(string(n))
	if len(data.Id) != 0 {
		return_id = GenerateShortUrl()
		return return_id
	}
	return_id = string(n)
	return return_id
}

// Check time expired
func ExpireData(date string) bool {
	var expired bool = false
	layout := "2006-01-02T03:04Z"
	expired_time, _ := time.Parse(layout, date)
	now := time.Now()
	expired = expired_time.Before(now)
	return expired
}
