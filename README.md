# Shorten URL

## 產生短網址
利用 crypto/rand package 使得 A-Z, a-z, 0-9 的字元隨機產生六位數的字串作為短網址的 ID
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
在 GenerateShortUrl( ) 設計中，經過隨機字串的產生後，需到資料庫進行比對，如果已經有相同的字串，便重新產生隨機字串，直到無重複為止。

此設計將遇到些問題，當資料庫已經存放上千筆資料過後，會導致將花費更多的時間找出尚未重複的隨機字串，且資料庫存放的空間將越來越少。如果可以，還可以再設計一些功能：當透過短網址 GET 一筆資料時，紀錄此筆資料透過短網址搜尋的點擊次數，有了此數據就能定期的將點擊次數為零或是最少的資料刪除，減少資料庫的負擔。

## 時間處理
在Go語言中，時間包提供了確定和查看時間的函數，在time package 中的 Parse 函數用於解析格式化的字符串，然後查找它形成的時間值，layout 通過以哪種方式顯示參考時間(即定義為 Mon Jan 2 15:04:05 -0700 MST 2006)來指定格式。

想將得知 expireAt 到期與否，可以使用套件提供的 time.now( ) 獲得當下時間，再使用 comp_time.Before(now) 檢查時間是不是有超過現在的時間
```
func ExpireData(date string) bool{
	var expired bool = false
	layout := "2006-01-02T03:04Z"
	expired_time, _ := time.Parse(layout, date)
	now := time.Now()
	expired = expired_time.Before(now)
	return expired
}
```

另外，目前經由短網址重新導向 URL 前會先檢查 expireAt 是否已經到期在進行動作，如果還有充裕的時間，能再將資料庫系統設計的更周全，當 expired date 已經到了期限時，會將該筆資料從資料庫中刪除，也可以解決上述短網址產生數量有限的問題。

## Gin
使用golang打造的web框架有很多種，例如 Beego、Echo 與 Gin，在使用 benchmark 的情況下 Echo 與 Gin 的效能都已經佔優勢了，兩個的框架寫法相似且簡潔，但 gin 內建了返回 html 文件的方法，因此更適合運用在 redirect URL(HTTP 302) 的短網址服務API。

促使我選擇 gin 來開發 golang API 是因為 gin 有以下主要的特點
* 基於原生的 net/http package 進行封裝
* 優秀的性能表現提供好的服務
* 使用極其快速的 http router
* http request 與 response 的驚喜包 - gin.Context


#### 匯入gin package
一開始要使用 gin 需要進行 import，因此我們就先將 package 進行 import
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
Golang 可以通過 database/sql package 實現了對 RDBMS 的使用，在 golang 中操作 mysql 資料庫比較簡單，package 本身也是使用 go 寫的，golang 原生有提供關於 sql 的抽象介面 database/sql，但後來有人利用他封裝了 go-sql-driver 支援 database/sql，我們會利用這個 package 進行實作。
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
透過 mysql 的 driver 建立的話，他有內建 Exec 的方法，可以直接執行原生的SQL指令，因此只要建立一個方法名為 CreateTable，然後把一開始連線建立好的 DB 當作參數傳入，之後再利用 Exec 的指令建立 Table

Table 中設計四欄位，分別為 url, id, shortUrl, expireAt 來存放 URL 與對應的 id 資訊，如下
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
經過 POST 所獲得的 url 與 expireAt 透過產生對應的網址id後，將 url, id, shortUrl, expireAt 要存入資料庫的值使用 Exec 的指令即可新增資料
```
_,err := DB.Exec("insert INTO urls(url, id, shortUrl, expireAt) values(?,?,?,?)",url, id, shortUrl, expireAt)
if err != nil{
    fmt.Printf("Create url ERROR:%v", err)
    return err
}
return nil
```
#### 查詢資料
Driver 他也提供了 Query 的語法可以供我們進行查詢，首先要定義搜尋回來的資料結構，URL 有 Long_URL, Id, Short_URL 與 ExpiredDate 四個參數，因此我們可以建立一個 struct 為
```
type URL struct{
	Long_URL string 
	Id string 
	Short_URL string 
	ExpiredDate string 
}
```
透過 Query 的方法執行 select 指令，支援將 Where 的值抽出來作為變數，我們這邊可分為兩種不一樣的查詢，分別為 QueryId( ) 與QueryUrl( )

QueryUrl( ) 為從原始的 URL 在資料庫進行查詢，而得到 struct 結構的資料
QueryId( ) 則是透過產生透過網址的id在資料庫中進行查詢，而得到整筆資料
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

## 短網址服務的 RESTful API 
<http://www.zoe-lin.me/static/index.html>  
先前已經有使用 Python Flask 開發過短網址服務，能將過將原本的網址縮短成精簡的短網址，當短網址被使用時，系統會先查出原本的網址，再以 URL 重新導向(HTTP 302)來將縮短後的位址重新導向到原來的 URL，並建立了一個 PostgreSQL 資料庫，儲存⻑網址與短網址之間相對應的資料。此服務也包括紀錄使用短網址搜尋的點擊次數、儲存搜尋短網址的使用者表頭資料，以計算出點擊率並且得知點擊者的時間分佈和族群以及前端的 HTML/CSS/JavaScript 製作，最後把已設計完成的 API 部署至 Heroku 雲平台上，將自己的網域轉向至 Heroku 的 API 上。

在時間充足的情況下，接下來的目標是希望能夠運用上述的設計開發 Golang Gin 的短網址服務，提供更完整的使用介面和體驗。