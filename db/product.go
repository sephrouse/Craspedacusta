package db

import (
	"fmt"
	"time"
)

type Product struct {
	Id           int       `xorm:"pk autoincr"`
	Father       string    `xorm:"varchar"`
	Title        string    `xorm:"varchar"`
	Url          string    `xorm:"varchar"`
	Rank         int       `xorm:"int"`
	Star         string    `xorm:"varchar"`
	Price        string    `xorm:"varchar"`
	Manufacturer string    `xorm:"varchar"`
	Parameters   string    `xorm:"varchar"`
	ImageUrl     string    `xorm:"varchar"`
	CreateTime   time.Time `xorm:"created"`
}

func InsertProduct(products []Product, tableName string) error {
	affected, err := engine.Table(tableName).Insert(&products)
	if err != nil {
		fmt.Println("db.InsertProduct failed. ", err)
		return err
	}

	fmt.Println("db.InsertProduct succeeded. ", affected, " rows have been inserted.")
	return nil
}
