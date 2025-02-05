package main

import (
	"fmt"
	"github.com/ensarkovankaya/go-messagingapp/common"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(postgres.Open(common.Cnf.DatabaseURI), &gorm.Config{
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}
}
