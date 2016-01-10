package main

import(
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var db *sql.DB  //全局变量

func init(){
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test")
	if db == nil {
		panic("open db error")
	}

	fmt.Println("open db ok")
	db.SetMaxOpenConns(100) //最多100个mysql数据库连接
	db.SetMaxIdelConns(10)  //最多10个空闲连接
}

func ConnectDB(){

	fmt.Println("test db ping")
	err := db.Ping()
	if err != nil {
		fmt.Println("database can not connect")
		return
	}
	fmt.Println("db connect ok")

	fmt.Println("test sql execute")
	rows, err := db.Query("select * from user limit 1")
	defer rows.Close()
	
	if err != nil {
		panic(err)
	}


}

func main(){
	ConnectDB()
}