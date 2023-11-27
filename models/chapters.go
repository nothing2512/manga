package models

type Chapter struct {
	ID      int    `json:"id" gorm:"primary_key"`
	MangaId int    `json:"mangaId"`
	Link    string `json:"link"`
}

func (*Chapter) TableName() string {
	return "chapters"
}
