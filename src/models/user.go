package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;not null"`
	FirstName    string    `db:"first_name" gorm:"not null"`
	LastName     string    `db:"last_name" gorm:"not null"`
	Email        string    `db:"email" gorm:"not null;unique"`
	Password     string    `db:"password" gorm:"not null"`
	IsAmbassador bool      `db:"is_ambassador" gorm:"default:false"`
}

func (u *User) SetPassword(p string) error {
	password, err := bcrypt.GenerateFromPassword([]byte(p), 12)
	if err != nil {
		return err
	}

	u.Password = string(password)

	return nil
}

func (u *User) ComparePassword(p string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
}
