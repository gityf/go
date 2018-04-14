package main

import (
	"database/sql"
	_ "github.com/mysql"
	"fmt"
)

func main() {
	fmt.Println("hi mysql")
	db, err := sql.Open("mysql","root:123456@tcp(192.10.1.120:3306)/meetme")
	if err != nil {
		fmt.Printf("errorcode %s", err.Error())
	}
	defer db.Close()
	stmtQry, err := db.Prepare("select id, admin, password from user where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtQry.Close()
	var id int
	var admin,passwd string
	err = stmtQry.QueryRow(20).Scan(&id, &admin,&passwd)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("id:%d,admin:%s,pass:%s\n", id, admin, passwd)	
	db.Exec("create table test(id int, name varchar(10))")
	db.Exec("insert into test(id, name) values(1,'w1'),(2,'w2')")
	rows, err := db.Query("select id, name from test")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	var name string
	for rows.Next() {
		rows.Scan(&id, &name)
		fmt.Printf("test.id:%d, name:%s", id, name)
	}
	
	stmtIns, err := db.Prepare("insert into test values(?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()
	for i:=10; i < 15; i++ {
		name = fmt.Sprintf("wangyf_%d", i)
		stmtIns.Exec(i, name)
	}
}