package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dns := "develop:password@tcp(db:3306)/mydb"
	db, err := sql.Open("mysql", dns)
	if err != nil{
		log.Fatalf("接続失敗: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ping失敗： %v", err)
	}
	fmt.Println("MySQLに接続成功!")
}