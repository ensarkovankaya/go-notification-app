package main

import (
	"fmt"
	"github.com/ensarkovankaya/go-notification-app/repositories"
	"go.uber.org/zap"
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

func closeDB() {
	sql, err := DB.DB()
	if err != nil {
		zap.L().Error("failed to get database connection", zap.Error(err))
	}
	if err = sql.Close(); err != nil {
		zap.L().Error("failed to close database connection", zap.Error(err))
	}
}
