# Shorten URL


## 產生短網址
利用crypto/rand package 使得A~Z, a~z, 0~9的字元隨機產生六位數的字串作為短網址的ID
```
const base = 36
size := big.NewInt(base)
n := make([]byte, 6)
for i, _ := range n {
    c, _ := rand.Int(rand.Reader, size)
    n[i] = strconv.FormatInt(c.Int64(), base)[0]
}
return string(n)
```
在GenerateShortUrl()設計中，經過隨機字串的產生後，需到資料庫進行比對，如果已經有相同的字串，便重新產生隨機字串，直到無重複為止。

此設計將遇到些問題，當資料庫已經存放上千筆資料過後，會導致將花費更多的時間找出尚未重複的隨機字串，且資料庫存放的空間將越來越少。如果可以，還可以再設計一項功能，當expired date已經到了期限時，會將該筆資料在資料庫中刪除。

## Gin
使用golang打造的web框架有很多種，例如Beego、Echo與Gin，在使用benchmark的情況下Echo與Gin的效能都已經佔優勢了，兩個的框架寫法相似且簡潔，但gin內建了返回html文件的方法，因此更適合運用在redirect URL(HTTP 302)的短網址服務API。

促使我選擇gin來開發golang API是因為gin有以下主要的特點
* 基於原生的net/http package進行封裝
* 優秀的性能表現提供好的服務
* 使用極其快速的http router
* http request與response的驚喜包 - gin.Context


#### 匯入gin package
一開始要使用gin需要進行import，因此我們就先將package進行import
```
import (
    "github.com/gin-gonic/gin"
    "net/http"
)
```
#### 設定 http routing
```
router.POST("/api/v1/urls", api_service.PostUrl)
router.GET("/:url_id", api_service.GetUrl)
```
#### 啟動Gin server
```
server.Run(":8000")
```

## MySQL
Golang可以通過database/sql package實現了對RDBMS的使用，在golang中操作mysql資料庫比較簡單，package本身也是使用go寫的，golang原生有提供關於sql的抽象介面database/sql，但後來有人利用他封裝了go-sql-driver支援database/sql，我們會利用這個package進行練習。
#### 連線
首先我們要先匯入 database/sql 與 github.com/go-sql-driver/mysql，mysql driver 前面要加上 _
```
    import (
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
    )
```
```
	conn := "<user_name>:<password>@<network>(<server>:<port>)/<database>"
	db, err := sql.Open("mysql", conn)
```
#### 建立Table
透過mysql的driver建立的話，他有內建Exec的方法，可以直接執行原生的SQL指令，因此只要建立一個方法名為CreateTable，然後把一開始連線建立好的DB當作參數傳入，之後再利用Exec的指令建立Table

Table中設計四欄位，分別為url, id, shortUrl, expireAt來存放URL與對應的id資訊，如下
```
func CreateTable() error {
	sql := `CREATE TABLE IF NOT EXISTS urls(
        url VARCHAR(200),
        id VARCHAR(6) PRIMARY KEY,
        shortUrl VARCHAR(30),
        expireAt VARCHAR(20)
	); `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("Create table ERROR:", err)
		return err
	}
	return nil
}
```
#### 新增資料
經過POST所獲得的url與expireAt透過產生對應的網址id後，將url, id, shortUrl, expireAt要存入資料庫的值使用Exec的指令即可新增資料
```
_,err := DB.Exec("insert INTO urls(url, id, shortUrl, expireAt) values(?,?,?,?)",url, id, shortUrl, expireAt)
if err != nil{
    fmt.Printf("Create url ERROR:%v", err)
    return err
}
return nil
```
#### 查詢資料
Driver他也提供了Query的語法可以供我們進行查詢，首先要定義搜尋回來的資料結構，URL有Long_URL, Id, Short_URL 與 ExpiredDate四個參數，因此我們可以建立一個struct為
```
type URL struct{
	Long_URL string 
	Id string 
	Short_URL string 
	ExpiredDate string 
}
```
透過Query的方法執行select指令，支援將Where的值抽出來作為變數，我們這邊可分為兩種不一樣的查詢，分別為QueryId()與QueryUrl()

QueryUrl()為從原始的URL在資料庫進行查詢，而得到struct結構的資料
QueryId()則是透過產生透過網址的id在資料庫中進行查詢，而得到整筆資料
```
// Query by long url
func QueryUrl(long_url string) pojo.URL{
	url := new(pojo.URL)
	mapping := DB.QueryRow("select * from urls where url=?", long_url)
	mapping.Scan(&url.Long_URL, &url.Id, &url.Short_URL, &url.ExpiredDate);
	return *url
}

// Query by short url
func QueryId(short_url_id string) pojo.URL{
	url := new(pojo.URL)
	mapping := DB.QueryRow("select * from urls where id=?", short_url_id)
	mapping.Scan(&url.Long_URL, &url.Id, &url.Short_URL, &url.ExpiredDate);
	return *url
}
```
