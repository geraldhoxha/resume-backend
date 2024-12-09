package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	// "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db *gorm.DB
)

func GetDB() *gorm.DB {
	return db
}

func init() {
	connectDatabase()
}

func connectDatabase() {
	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_DATABASE")
		port     = os.Getenv("DB_PORT")
	)
	psqlDns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
	var err error

	// db, err = gorm.Open(mysql.Open(databaseConfig), initConfig())
	db, err = gorm.Open(postgres.Open(psqlDns), initConfig())
	if err != nil {
		fmt.Printf("%s. Err: %s", psqlDns, err.Error())
		panic("Fail To Connect Database")
	}
}

func initConfig() *gorm.Config {
	return &gorm.Config{
		Logger:         initLog(),
		NamingStrategy: initNamingStrategy(),
	}
}

func initLog() logger.Interface {
	f, _ := os.Create("gorm.log")
	newLogger := logger.New(log.New(io.MultiWriter(f), "\r\n", log.LstdFlags), logger.Config{
		Colorful:      true,
		LogLevel:      logger.Info,
		SlowThreshold: time.Second,
	})
	return newLogger
}

func initNamingStrategy() *schema.NamingStrategy {
	return &schema.NamingStrategy{
		SingularTable: false,
		TablePrefix:   "",
	}
}
