package databases

import (
	"github.com/wuttinanhi/code-judge-system/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSQLiteDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect sqlite database")
	}

	StartMigration(db)

	return db
}

func NewTempSQLiteDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// disable logger
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to create temp sqlite database")
	}

	StartMigration(db)

	return db
}

func StartMigration(db *gorm.DB) {
	db.AutoMigrate(
		&entities.User{},
		&entities.Challenge{},
		&entities.ChallengeTestcase{},
		&entities.Submission{},
		&entities.SubmissionTestcase{},
	)
}
