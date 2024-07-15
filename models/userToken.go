package models

import (
	"gorm.io/gorm"
)

type UserToken struct {
	gorm.Model
	UserID uint
	Token  string
}
