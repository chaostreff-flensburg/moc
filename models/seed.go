package models

import (
	"github.com/jinzhu/gorm"
)

// Seed generate a testentry
func Seed(db *gorm.DB) *Message {
	message := NewMessage("Hallo ich bin eine Test Nachricht")

	db.Create(message)

	return message
}
