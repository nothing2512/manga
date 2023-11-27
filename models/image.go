package models

type Image struct {
	ID        int    `json:"id" gorm:"primary_key"`
	ChapterId int    `json:"chapterId"`
	Link      string `json:"link"`
}

func (*Image) TableName() string {
	return "images"
}
