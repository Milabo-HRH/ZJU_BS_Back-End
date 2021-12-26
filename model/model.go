package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Mail      string `gorm:"type:varchar(30);not null;unique"`
	Password  string `gorm:"size:255;not null"`
	Privilege string `gorm:"type:varchar(10);not null"`
}

type Picture struct {
	gorm.Model
	UploaderID int    `gorm:"type:int unsigned;notnull"`
	FileName   string `gorm:"type:varchar(40);notnull"`
}

type Task struct {
	gorm.Model
	TaskName     string `gorm:"type:varchar(40);notnull"`
	PublisherID  int    `gorm:"type:int unsigned;notnull"`
	PictureID    int    `gorm:"type:int unsigned;notnull"`
	Tags         string `gorm:"type:varchar(40);"`
	Reviewed     bool   `gorm:"type:bool"`
	ReviewUserID int    `gorm:"type:int unsigned;notnull"`
}
