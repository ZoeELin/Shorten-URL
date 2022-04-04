package db_service

import (
    "fmt"
    "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golangAPI/pojo"
)

var DB *sql.DB

func DatabaseConnect(){
	conn := "root:Zozo0411@tcp(127.0.0.1:3306)/shorten_url"
	db, err := sql.Open("mysql", conn)
	DB = db
	
	if err != nil {
		fmt.Println("開啟 MySQL 連線發生錯誤，原因為：", err)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println("資料庫連線錯誤，原因為：", err.Error())
		return
	}
	// defer db.Close()
}

func CloseDatabase(){
	defer DB.Close()
}

func CreateTable() error {
	sql := `CREATE TABLE IF NOT EXISTS urls(
        url VARCHAR(64),
        id VARCHAR(6) PRIMARY KEY,
        shortUrl VARCHAR(30),
        expireAt VARCHAR(20)
	); `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("建立 Table 發生錯誤:", err)
		return err
	}
	fmt.Println("建立 Table 成功！")
	return nil
}

// Add into database
func InsertURL(url, id, shortUrl, expireAt string) error{
	_,err := DB.Exec("insert INTO urls(url, id, shortUrl, expireAt) values(?,?,?,?)",url, id, shortUrl, expireAt)
	if err != nil{
		fmt.Printf("建立使用者失敗，原因是：%v", err)
		return err
	}
	fmt.Println("建立使用者成功！")
	return nil
}

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

// Query all
func QueryAll(){
	url := new(pojo.URL)
	mapping := DB.QueryRow("select * from urls")
	mapping.Scan(&url.Long_URL, &url.Id, &url.Short_URL, &url.ExpiredDate);
	fmt.Println(*url)
	// return *url
}