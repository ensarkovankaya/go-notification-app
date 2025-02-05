package main

import (
	"fmt"
	"github.com/ensarkovankaya/go-notification-app/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(postgres.Open(Cnf.DatabaseURI), &gorm.Config{
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}
	if err = DB.AutoMigrate(&repositories.Message{}); err != nil {
		panic(fmt.Errorf("failed to migrate database: %w", err))
	}
}
