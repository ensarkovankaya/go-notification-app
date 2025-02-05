package main

import (
	"fmt"
	"github.com/Netflix/go-env"
	"time"
)

var Cnf Config

type Config struct {
	// Http Server Configurations
	Port         string        `env:"PORT,default=9098"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT,default=60s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT,default=60s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT,default=60s"`
	// Database Configurations
	DatabaseURI string        `env:"DB_URI"`
	WebhookID   string        `env:"WEBHOOK_ID"`
	RedisURI    string        `env:"REDIS_URI"`
	CronTTL     time.Duration `env:"CRON_TTL,default=2m"`
	LogLevel    string        `env:"LOG_LEVEL,default=ERROR"`
}

func init() {
	_, err := env.UnmarshalFromEnviron(&Cnf)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
}
