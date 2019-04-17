package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	v "gopkg.in/go-playground/validator.v9"

	"github.com/chaostreff-flensburg/moc/validator"
)

type MessageRequest struct {
	Text string `json:"message" validate:"required,min=3,max=160"`
}

type Message struct {
	MessageRequest

	ID string `gorm:"type:uuid; primary_key" json:"id"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NewMessage create new message
func NewMessage(text string) *Message {
	return &Message{
		MessageRequest: MessageRequest{
			Text: text,
		},
	}
}

// BeforeCreate will create a uuid right before creating
func (m *Message) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New().String())

	return nil
}

// Validate request by annotations
func (r *MessageRequest) Validate() *map[string]string {
	validate := validator.NewValidator()

	err := validate.Struct(r)
	if err != nil {
		errors := map[string]string{}

		for _, err := range err.(v.ValidationErrors) {
			errors[err.Field()] = err.ActualTag()
		}

		return &errors
	}

	return nil
}
