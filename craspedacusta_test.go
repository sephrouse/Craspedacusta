package main

import (
	"fmt"
	"strings"
	"testing"
)

func Test1(t *testing.T) {
	var tempTableName string

	for i := range baseCatagory {
		tempTableName = baseCatagory[i]
		tempTableName = strings.TrimSpace(tempTableName)
		tempTableName = strings.Replace(tempTableName, " ", "", -1)
		tempTableName = strings.Replace(tempTableName, "&", "", -1)

		fmt.Println(tempTableName)
	}
}
