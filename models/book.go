package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	ID         int
	BookName   string
	Author     string
	OriginalID string
}
