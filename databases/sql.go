package databases

import (
	"github.com/wuttinanhi/code-judge-system/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLiteDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect sqlite database")
	}

	// Migrate the schema
	db.AutoMigrate(&entities.User{})

	return db
}

func NewTempSQLiteDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to create temp sqlite database")
	}

	// Migrate the schema
	db.AutoMigrate(&entities.User{})

	return db
}
