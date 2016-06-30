package db

import (
	"fmt"
	"time"
)

type Catagory struct {
	Id         int       `xorm:"pk autoincr"`
	Name       string    `xorm:"varchar"`
	Level      int       `xorm:"int"`
	FatherName string    `xorm:"varchar"`
	Url        string    `xorm:"varchar"`
	CreateTime time.Time `xorm:"created"`
}

func InsertCatagory(catagories []Catagory) error {
	affected, err := engine.Insert(&catagories)
	if err != nil {
		fmt.Println("db.InsertCatagory failed. ", err)
		return err
	}

	fmt.Println("db.InsertCatagory succeeded. ", affected, " rows have been inserted.")
	return nil
}
