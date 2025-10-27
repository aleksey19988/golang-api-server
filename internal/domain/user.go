package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `json:"name"  gorm:"not null;size(255)"`
	Age   uint   `json:"age"   gorm:"not null;default:0"`
	Email string `json:"email" gorm:"uniqueIndex;type:varchar(255)"`
}
