package config

import (
	"fmt"
	"khaf-dev/model"
	"log"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type users model.Users

func init() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=kahfdev password=welcome1 dbname=kahfdev sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Db Connected....")
	}
	db.AutoMigrate(&users{})
}
