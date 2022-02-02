package models

import "gorm.io/gorm"

type BookChar struct {
	gorm.Model
	BookID   int
	BookChar string
	BookVal  string
}
