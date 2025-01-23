package db

import (
    "fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
}

func NewDB() *gorm.DB {
	if os.Getenv("GO_ENV") == "dev" {
		err := gotoenv.Load{}
		if err!= nil {
            log.Fatal(err)
        }
	}
}
