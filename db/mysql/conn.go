package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Printf("Failed to connect to mysql,err:" + err.Error())
		os.Exit(1)
	}
}

//DBConn:返回数据库连接对象
func DBConn() *sql.DB {
	return db
}
