package handlers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/sdejesusp/petsitter/database"
	"github.com/sdejesusp/petsitter/models"
)

// USER error codes
// Range 1000 - 1999
const (
	TokenBodyParse         = 1000
	TokenInvalidEmail      = 1001
	TokenBadCredencial     = 1002
	UserCannotCreate       = 1003
	UserUnauthorizeRequest = 1004
	UserInvalidId          = 1005
	UserNotFound           = 1006
	ModifyBodyParse        = 1007
)

type User struct {
	ID       uint64   `json:"id"`
	FullName string   `json:"fullName"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles" gorm:"serializer:json"`
}

type AllUsers struct {
	TotalCount  int    `json:"totalCount"`
	RecordIndex int    `json:"recordIndex"`
	Users       []User `json:"users" gorm:"serializer:json"`
}

const JWTSECRETENV = "JWTSECRET"

func ReadEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error reading jwt secret string")
	}

	return os.Getenv(key)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUserResponse(userModel models.User) User {
	return User{ID: userModel.ID, FullName: userModel.FullName, Email: userModel.Email, Roles: userModel.Roles}
}

func FindUser(id uint64, user *models.User) error {
	database.DB.Db.Find(&user, "id=?", id)
	if user.ID == 0 {
		return errors.New("user not found")
	}
	return nil
}

func FindUserWithEmail(email string, user *models.User) error {
	database.DB.Db.Find(&user, "email=?", email)
	if user.ID == 0 {
		return errors.New("user not found")
	}
	return nil
}

func GetToken(c *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":   "Cannot parse json",
			"errorCode": TokenBodyParse,
		})
	}

	body.Email = strings.ToLower(body.Email)

	var user models.User
	if err := FindUserWithEmail(body.Email, &user); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message":   "Bad Credentials",
			"errorCode": TokenInvalidEmail,
		})
	}

	if body.Email != user.Email || !CheckPasswordHash(body.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message":   "Bad Credentials",
			"errorCode": TokenBadCredencial,
		})
	}

	// Create the Claims
	expiry := time.Now().Add(time.Hour * 1).Unix()
	claims := jwt.MapClaims{
		"name":  user.FullName,
		"sub":   strconv.FormatUint(user.ID, 10),
		"admin": false,
		"exp":   expiry,
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := ReadEnvVariable(JWTSECRETENV)
	fmt.Println("jwtSecret", jwtSecret)

	// Generate encoded token and send it as response
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{
		"token":     t,
		"tokenType": "Bearer",
		"exp":       expiry,
	})

}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	hashPassword, err := HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":   "User cannot be create",
			"errorCode": UserCannotCreate,
		})
	}

	user.Email = strings.ToLower(user.Email)
	user.Password = hashPassword

	if database.DB.Db.Create(&user); user.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":   "User cannot be create",
			"errorCode": UserCannotCreate,
		})
	}

	response := CreateUserResponse(user)

	return c.Status(fiber.StatusCreated).JSON(response)
}

func GetUserProfile(c *fiber.Ctx) error {
	jwtUser := c.Locals("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)

	id, _ := strconv.ParseUint(sub, 10, 64)

	var user models.User
	if err := FindUser(id, &user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	response := CreateUserResponse(user)

	return c.Status(fiber.StatusOK).JSON(response)

}

func GetUsers(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	admin := claims["admin"].(bool)

	if !admin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message":   "Unauthorized request",
			"errorCode": UserUnauthorizeRequest,
		})
	}

	users := []models.User{}

	database.DB.Db.Find(&users)
	var response AllUsers

	responseUsers := []User{}
	for _, user := range users {
		responseUser := CreateUserResponse(user)
		responseUsers = append(responseUsers, responseUser)
	}

	response.Users = responseUsers
	response.TotalCount = len(responseUsers)
	response.RecordIndex = 20

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetUserWithId(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id < 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message":   "Please ensure that id is a positive integer",
			"errorCode": UserInvalidId,
		})
	}

	userId := uint64(id)

	var user models.User
	if err := FindUser(userId, &user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	response := CreateUserResponse(user)

	return c.Status(fiber.StatusOK).JSON(response)
}

func ModifyUserWithId(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":   err.Error(),
			"errorCode": UserInvalidId,
		})
	}

	userId := uint64(id)

	var user models.User
	if err := FindUser(userId, &user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message":   err.Error(),
			"errorCode": UserNotFound,
		})
	}

	type UpdateFields struct {
		FullName string   `json:"full_name"`
		Email    string   `json:"email"`
		Roles    []string `json:"roles" gorm:"serializer:json"`
	}

	var updateData UpdateFields

	if c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":   err.Error(),
			"errorCode": ModifyBodyParse,
		})
	}

	user.FullName = updateData.FullName
	user.Email = updateData.Email
	user.Roles = updateData.Roles

	database.DB.Db.Save(&user)
	response := CreateUserResponse(user)

	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteUserWithId(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":   err.Error(),
			"errorCode": UserInvalidId,
		})
	}

	userId := uint64(id)

	var user models.User
	if err := FindUser(userId, &user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	if err := database.DB.Db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	return c.Status(fiber.StatusNoContent).SendString("")
}
