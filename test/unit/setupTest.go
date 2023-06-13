package unit

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getMockDatabase() (sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return sql.DB{}, nil, err
	}

	gormDb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		return sql.DB{}, nil, err
	}

	mockDb, err := gormDb.DB()
	if err != nil {
		return sql.DB{}, nil, err
	}

	return *mockDb, mock, nil
}
