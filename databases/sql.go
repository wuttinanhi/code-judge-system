package databases

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/mysql"

	"github.com/spf13/viper"
	"github.com/wuttinanhi/code-judge-system/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartMigration(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.ChallengeTestcase{},
		&entities.Challenge{},
		&entities.SubmissionTestcase{},
		&entities.Submission{},
		&entities.User{},
	)
}

func NewSQLiteDatabase() *gorm.DB {
	cwd, err := os.Getwd()
	if err != nil {
		panic("failed to get current working directory")
	}

	sqlitepath := filepath.Join(cwd, "test.db")

	db, err := gorm.Open(sqlite.Open(sqlitepath), &gorm.Config{})
	if err != nil {
		panic("failed to connect sqlite database")
	}

	StartMigration(db)

	return db
}

func NewTempSQLiteDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		// disable logger
		// Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to create temp sqlite database")
	}

	StartMigration(db)

	return db
}

func NewMySQLDatabase() *gorm.DB {
	DB_HOST := viper.GetString("DB_HOST")
	DB_PORT := viper.GetString("DB_PORT")
	DB_USER := viper.GetString("DB_USER")
	DB_PASSWORD := viper.GetString("DB_PASSWORD")
	DB_NAME := viper.GetString("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DB_USER,
		DB_PASSWORD,
		DB_HOST,
		DB_PORT,
		DB_NAME,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect mysql database")
	}

	err = StartMigration(db)
	if err != nil {
		panic("failed to migrate mysql database")
	}

	return db
}
