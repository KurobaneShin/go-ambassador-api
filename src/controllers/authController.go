package controllers

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"ambassador/src/database"
	"ambassador/src/models"
)

type RegisterParams struct {
	FirstName       string `json:"first_name" valid:"type(string),required"`
	LastName        string `json:"last_name" valid:"type(string),required"`
	Email           string `json:"email" valid:"type(string),email,required"`
	Password        string `json:"password" valid:"type(string),required"`
	PasswordConfirm string `json:"password_confirm" valid:"type(string),required"`
}

type LoginParams struct {
	Email    string `json:"email" valid:"type(string),email,required"`
	Password string `json:"password" valid:"type(string),required"`
}

func Register(ctx *fiber.Ctx) error {
	var data map[string]string

	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	_, err := govalidator.ValidateStruct(RegisterParams{
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

	user := models.User{
		Id:           uuid.New(),
		FirstName:    data["first_name"],
		LastName:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: false,
	}

	err = user.SetPassword(data["password"])
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"message": "Could not create user",
		})
	}

	database.DB.Create(&user)

	return ctx.JSON(user)
}

func Login(ctx *fiber.Ctx) error {
	var data map[string]string

	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	_, err := govalidator.ValidateStruct(LoginParams{
		Email:    data["email"],
		Password: data["password"],
	})
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user := models.User{}

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == uuid.Nil {
		ctx.Status(fiber.StatusNotFound)

		return ctx.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		ctx.Status(fiber.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	payload := jwt.StandardClaims{
		Subject:   user.Id.String(),
		ExpiresAt: jwt.NewTime(15000),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)

		return ctx.JSON(fiber.Map{
			"message": "Could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Value:    signedToken,
	}

	ctx.Cookie(&cookie)

	return ctx.JSON(fiber.Map{
		"message": "Success",
		"token":   signedToken,
	})
}
