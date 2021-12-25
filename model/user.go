package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Mail      string `gorm:"type:varchar(30);not null;unique"`
	Password  string `gorm:"size:255;not null"`
	Privilege string `gorm:"type:varchar(2)"`
}
