package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model

	Blogtitle       string   `gorm:"size:255;uniqueIndex" json:"blogtitle" validate:"required,min=3,max=50"`
	Blogsubtitle    string   `gorm:"size:255;uniqueIndex" json:"blogsubtitle" validate:"required,min=3,max=50"`
	Blogimage       string   `gorm:"size:255;uniqueIndex" json:"blogimage" validate:"required,min=3,max=50"`
	Blogdescription string   `gorm:"size:255;uniqueIndex" json:"blogdescription" validate:"required,min=3,max=50000"`
	CategoryID      uint     `json:"categoryID"`
	Category        Category `gorm:"foreignKey:CategoryID"`
	UserID          uint     `json:"userID"`
	User            User     `gorm:"foreignKey:UserID"`
}

func (u *Blog) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
