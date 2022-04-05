# Shorten URL
## Gin
Gin是一套使用golang打造的web框架，gin主要的特點如下
* 優秀的性能表現提供好的服務
* 使用極其快速的httprouter
* http request與response的驚喜包 - gin.Context

#### 匯入gin package
一開始要使用gin需要進行import，因此我們就先將package進行import
```
import (
    "github.com/gin-gonic/gin"
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
透過程式語言操作資料庫，最常見的方法就是使用driver，golang原生有提供關於sql的抽象介面database/sql，但後來有人利用他封裝了 mysql的driver - go-sql-driver，我們會利用這個package進行練習。
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
透過mysql的driver建立的話，他有內建Exec的方法，可以直接執行原生的SQL指令，因此只要建立一個方法名為CreateTable，然後把一開始連線建立好的DB當作參數傳入，之後再利用Exec的指令建立Table即可
```
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
```
#### 新增資料
一樣透過 Exec 的指令即可
```
_,err := DB.Exec("insert INTO urls(url, id, shortUrl, expireAt) values(?,?,?,?)",url, id, shortUrl, expireAt)
if err != nil{
    fmt.Printf("Create url ERROR:%v", err)
    return err
}
return nil
```
#### 查詢資料
使用driver的時候，他有提供Query的語法可以供我們進行查詢，首先要定義搜尋回來的資料結構，URL有Long_URL, Id, Short_URL 與 ExpiredDate四個參數，因此我們可以建立一個struct為
```
type URL struct{
	Long_URL string 
	Id string 
	Short_URL string 
	ExpiredDate string 
}
```
透過 Query的方法執行 select 指令，Query 的方法支援將Where的值抽出來作為變數，我們這邊可分為兩種不一樣的查詢，分別為QueryId()與QueryUrl()
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
