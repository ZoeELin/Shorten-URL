package database

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func DatabaseConnect(){
	conn := "root:Zozo0411@tcp(127.0.0.1:3306)/shorten_url"
	db, err := sql.Open("mysql", conn)
	
	if err != nil {
		fmt.Println("開啟 MySQL 連線發生錯誤，原因為：", err)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println("資料庫連線錯誤，原因為：", err.Error())
		return
	}
	defer db.Close()
}