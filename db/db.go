package db

import (
	"fmt"

	"hamal/models"

	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func InitDB() {
	sqlc := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		"root", "root", "127.0.0.1", "3306", "hamal")
	var err error
	db, err = gorm.Open("mysql", sqlc)
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.User{}).Error; err != nil {
		log.Println(err.Error())
	}
	db.AutoMigrate(&models.Currency{})
	db.AutoMigrate(&models.Strategy{})
	db.AutoMigrate(&models.StrategyInstance{})
	db.AutoMigrate(&models.Symbol{})
	db.AutoMigrate(&models.DailyStat{})
	db.AutoMigrate(&models.StrategyInstanceRecard{})
	db.AutoMigrate(&models.Exchange{})
	db.AutoMigrate(&models.UserExchange{})

	fmt.Println("db init success.")
}

func OffsetLimit(page int, limit int) (int, int) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return (page - 1) * limit, limit
}
