package models

import "time"

type Manga struct {
	ID          int       `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Link        string    `json:"link"`
	Index       int       `json:"index"`
	LastUpdated time.Time `json:"lastUpdated"`
	LastRead    int       `json:"lastRead"`
}

func (*Manga) TableName() string {
	return "mangas"
}
