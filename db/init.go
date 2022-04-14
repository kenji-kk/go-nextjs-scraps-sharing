package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB
var err error 

func DbConnect() {
	Db, err = sql.Open("mysql", "go_grpc:password@tcp(mysql:3306)/go_database?charset=utf8&parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("db接続完了")

	cmdU := `CREATE TABLE IF NOT EXISTS users (
		id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE,
		hashedpassword LONGBLOB NOT NULL,
		salt LONGBLOB NOT NULL
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		`

	_, err = Db.Exec(cmdU)
	count := 0
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 180 {
				fmt.Println("")
				panic(err)
			}
			_, err = Db.Exec(cmdU)
		}
	}
	fmt.Println("ユーザテーブル作成成功")

	fmt.Println("テーブル作成成功")
}
