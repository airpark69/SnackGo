package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	UserName string
	Message  string
}
