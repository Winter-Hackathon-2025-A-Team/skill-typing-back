package db

import (
    "fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	if os.Getenv("GO_ENV") == "dev" {
		err := gotoenv.Load{}
		if err!= nil {
            log.Fatal(err)
        }
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PW"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_HOST"),
        os.Getenv("MYSQL_DB"))
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
        log.Fatal(err)
    }
	fmt.Println("Connected")
	return db
}
