package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Mail      string `gorm:"type:varchar(30);not null;unique" json:"Mail"`
	Password  string `gorm:"size:255;not null" json:"Password"`
	Privilege string `gorm:"type:varchar(10);not null"`
}

type Assignment struct {
	gorm.Model
	UploaderID uint   `gorm:"type:int unsigned;notnull"`
	Filename   string `gorm:"type:varchar(100);notnull"`
	Annotated  bool   `gorm:"type:bool"`
	Reviewed   bool   `gorm:"type:bool"`
	Tags       string `gorm:"type:varchar(40);"`
}

type Annotation struct {
	gorm.Model
	UploaderID   uint   `gorm:"type:int unsigned;notnull"`
	AssignmentID uint   `gorm:"type:int unsigned;notnull"`
	Tags         string `gorm:"type:varchar(40);"`
	Reviewed     bool   `gorm:"type:bool"`
	ReviewUserID uint   `gorm:"type:int unsigned;notnull"`
}
