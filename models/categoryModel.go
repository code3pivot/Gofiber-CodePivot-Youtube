package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model

	Categoryname string `gorm:"size:255;uniqueIndex" json:"categoryname" validate:"required,min=3,max=50"`
}

func (u *Category) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
