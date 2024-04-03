package models

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;not null"`
	FirstName    string    `db:"first_name" gorm:"not null"`
	LastName     string    `db:"last_name" gorm:"not null"`
	Email        string    `db:"email" gorm:"not null;unique"`
	Password     string    `db:"password" gorm:"not null"`
	IsAmbassador bool      `db:"is_ambassador" gorm:"default:false"`
}
