package database

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	entity "social-app/internal/entity"
)

func InitMysql() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("database_user"),
		os.Getenv("database_pass"),
		os.Getenv("database_host"),
		os.Getenv("database_port"),
		os.Getenv("database_name"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Errorf("Failed to connect to database: %s", err.Error())
		return nil, err
	}

	tx := db.Begin()

	if err := tx.AutoMigrate(
		&entity.User{},
		&entity.Post{},
		&entity.UserFollower{},
		&entity.UserLikedPost{},
		&entity.UserViewedPost{},
	); err != nil {
		log.Errorf("Failed to run database migration: %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	log.Info("Database migration successful")

	return db, nil
}
