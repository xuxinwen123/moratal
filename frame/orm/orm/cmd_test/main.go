package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"orm"
	"orm/log"
)

func main() {
	engine, err := orm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local&parseTime=true", "root", "123456", "localhost", 3306, "db_orm"))
	if err != nil {
		log.Error("failed to create conn")
		return
	}
	defer engine.Close()
	s := engine.NewSession()
	_, err = s.Raw("DROP table if exists User;").Exec()
	_, err = s.Raw("create table User(Name varchar(255),Age integer)").Exec()
	result, err := s.Raw("insert into User(`Name`,`Age`) VALUES (?,?),(?,?)", "tom", 18, "jack", 25).Exec()
	rowsAffected, err := result.RowsAffected()
	fmt.Printf("exec succsee ,%d affected", rowsAffected)
}
