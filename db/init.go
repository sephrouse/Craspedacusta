package db

import (
	"fmt"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("postgres", "user=lakeman password='MorningAngel' host=localhost port=5432 dbname=lake sslmode=disable")
	if err != nil {
		fmt.Println("db.init: error. ", err)
		panic(err)
	}

	fmt.Println("db.init: DB Connected.")

	syncTable()
}

func syncTable() error {
	err := engine.Sync(new(Catagory))
	if err != nil {
		fmt.Println("db.syncTable: synchronize Catagory table error. ", err)
		return err
	}

	fmt.Println("db.syncTable: synchronization completed.")
	return nil
}

func Close() {
	engine.Close()
}

func CreateProductTable(tableName string) error {
	has, err := engine.IsTableExist(tableName)
	if err != nil {
		fmt.Println("db.CreateProductTable error. ", err)
		return err
	}

	if has == false {
		err = engine.Table(tableName).CreateTable(new(Product))
		if err != nil {
			fmt.Println("db.CreateProductTable error. ", err)
			return err
		}

		fmt.Println("db.CreateProductTable CreateTable successful.")
	} else {
		fmt.Println("db.CreateProductTable: ", tableName, " is existed.")
	}

	return nil
}
