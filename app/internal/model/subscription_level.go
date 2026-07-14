package model

type SubscriptionLevel struct {
	Model
	Name  string `json:"name" gorm:"not null;column:name"`
	Level int    `json:"level" gorm:"not null;uniqueIndex"`
	Limit int    `json:"limit" gorm:"not null"`
}
