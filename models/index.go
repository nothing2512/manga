package models

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s database=%s port=%s TimeZone=Asia/Jakarta",
		"localhost",
		"robet",
		"postgres",
		"manga",
		"manga",
		"5432",
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		&Chapter{},
		&Image{},
		&Manga{},
	)
	if err != nil {
		panic(err)
	}

	return db
}
