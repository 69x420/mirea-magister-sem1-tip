package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zaldk/mirea-sem1-tid-private/pr6/internal/models"
)

const DefaultDSN = "file:pz6_gorm.db?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=1"

func Connect() *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = DefaultDSN
	}

	d, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal("connect db: ", err)
	}

	sqlDB, err := d.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to SQLite via GORM")
	return d
}

func Migrate(d *gorm.DB) error {
	if err := d.AutoMigrate(&models.User{}, &models.Note{}, &models.Tag{}); err != nil {
		return err
	}
	log.Println("AutoMigrate is done")
	return nil
}
