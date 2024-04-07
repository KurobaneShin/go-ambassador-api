package controllers

import (
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"ambassador/src/database"
	"ambassador/src/models"
)

type RegisterPayload struct {
	FirstName       string `json:"first_name" valid:"type(string),required"`
	LastName        string `json:"last_name" valid:"type(string),required"`
	Email           string `json:"email" valid:"type(string),email,required"`
	Password        string `json:"password" valid:"type(string),required"`
	PasswordConfirm string `json:"password_confirm" valid:"type(string),required"`
}

func Register(ctx *fiber.Ctx) error {
	var data map[string]string

	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	_, err := govalidator.ValidateStruct(RegisterPayload{
		FirstName:       data["first_name"],
		LastName:        data["last_name"],
		Email:           data["email"],
		Password:        data["password"],
		PasswordConfirm: data["password_confirm"],
	})
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if data["password"] != data["password_confirm"] {
		ctx.Status(fiber.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"message": "Password do not match",
		})
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
	}

	user := models.User{
		Id:           uuid.New(),
		FirstName:    data["first_name"],
		LastName:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: false,
		Password:     string(password),
	}

	database.DB.Create(&user)

	return ctx.JSON(user)
}
