package config

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Configs struct {
	DbUsername string
	DbPassword string
	DbName     string
	DbHost     string
	DbSSL      string
}

func ReadConfig() *Configs {
	return &Configs{
		DbUsername: os.Getenv("DB_USERNAME"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		DbHost:     os.Getenv("DB_HOST"),
		DbSSL:      os.Getenv("DB_SSL"),
	}
}

func Connect(config *Configs) (*gorm.DB, error) {

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:1433?database=%s", config.DbUsername, config.DbPassword, config.DbHost, config.DbName)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
