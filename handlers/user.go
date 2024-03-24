package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sdejesusp/petsitter/database"
	"github.com/sdejesusp/petsitter/models"
)

type User struct {
	ID       uint     `json:"id"`
	FullName string   `json:"full_name"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles" gorm:"serializer:json"`
}

func CreateUserResponse(userModel models.User) User {
	return User{ID: userModel.ID, FullName: userModel.FullName, Email: userModel.Email, Roles: userModel.Roles}
}

func FindUser(id int, user *models.User) error {
	database.DB.Db.Find(&user, "id=?", id)
	if user.ID == 0 {
		return errors.New("user not found")
	}
	return nil
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	database.DB.Db.Create(&user)
	response := CreateUserResponse(user)

	return c.Status(fiber.StatusCreated).JSON(response)
}

func GetUserWithId(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON("please ensure that id is an integer")
	}

	var user models.User
	if err := FindUser(id, &user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	response := CreateUserResponse(user)

	return c.Status(fiber.StatusOK).JSON(response)
}

func ModifyUserWithId(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	var user models.User
	if err := FindUser(id, &user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	type UpdateFields struct {
		FullName string   `json:"full_name"`
		Email    string   `json:"email"`
		Roles    []string `json:"roles" gorm:"serializer:json"`
	}

	var updateData UpdateFields

	if c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	user.FullName = updateData.FullName
	user.Email = updateData.Email
	user.Roles = updateData.Roles

	database.DB.Db.Save(&user)
	response := CreateUserResponse(user)

	return c.Status(fiber.StatusOK).JSON(response)
}
