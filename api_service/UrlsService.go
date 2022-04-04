package api_service

import (
	"golangAPI/pojo"
	"golangAPI/db_service"
	"net/http"
	"github.com/gin-gonic/gin"
	"crypto/rand"
	"math/big"
	"strconv"
	"fmt"
)


// Get all urls
func FindAllUrls(c *gin.Context){
	db_service.QueryAll()
}

// POST LongURL and Expired data
func PostUrl(c *gin.Context){
	url := pojo.URL{}
	c.BindJSON(&url)

	data := db_service.QueryUrl(url.Long_URL)
	if len(data.Long_URL) == 0{
		url.Id = GenerateShortUrl()
		url.Short_URL = "http://localhost:8000/" + url.Id

		fmt.Println("ID: " + url.Id)
		fmt.Println("Shorten URL: " + url.Short_URL)

		db_service.InsertURL(url.Long_URL, url.Id, url.Short_URL, url.ExpiredDate)
	}else {
		url = data
	}

	c.JSON(http.StatusOK, gin.H{
		"id": url.Id,
		"shortUrl": url.Short_URL,
	})
}

// GET shortUrl and redirect to original URL
func GetUrl(c *gin.Context){
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request

	url_id := c.Param("url_id")
	data := db_service.QueryId(url_id)
	if len(data.Id) != 0{
		http.Redirect(w, r, data.Long_URL, 302)
	}else {
		c.JSON(http.StatusNotAcceptable, "The short URL is not exist.")
	}
}

// Create shorter URL
func GenerateShortUrl() string {
	const base = 36
	size := big.NewInt(base)
	n := make([]byte, 6)
	for i, _ := range n {
		c, _ := rand.Int(rand.Reader, size)
		n[i] = strconv.FormatInt(c.Int64(), base)[0]
	}
	return string(n)
}

// Check time expired?
// func ExpiredData(time string) bool{

// }

