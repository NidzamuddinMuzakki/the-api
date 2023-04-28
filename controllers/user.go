package controllers

import (
	"github.com/NidzamuddinMuzakki/the-api/configs"
	"github.com/NidzamuddinMuzakki/the-api/models"
	"github.com/NidzamuddinMuzakki/the-api/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsers(c *fiber.Ctx) error {
	var users []models.Profile

	configs.Database.Table("users").Find(&users)
	data := responses.UserResponse{
		Status:  200,
		Message: "success",
		Data:    users,
	}
	return c.Status(200).JSON(data)
}

func AddUser(c *fiber.Ctx) error {
	user := new(models.User)
	id := uuid.New()
	data := responses.UserResponse{
		Status: fiber.ErrBadRequest.Code,

		Data: nil,
	}
	if err := c.BodyParser(user); err != nil {
		data.Message = err.Error()
		return c.Status(fiber.ErrBadRequest.Code).JSON(data)
	}
	file, err := c.FormFile("profile_image")
	if err != nil {
		data.Message = err.Error()
		return c.Status(fiber.ErrBadRequest.Code).JSON(data)

	}
	user.ProfileImage = "uploads/" + id.String() + ".jpg"
	c.SaveFile(file, user.ProfileImage)
	// user.ProfileImage = "uploads/" + id.String() + ".jpg"

	// fmt.Println(user.ProfileImage, user.FirstName, user.Las)
	datas := configs.Database.First(&user, "username=?", user.Username)
	if datas.RowsAffected != 0 {
		data.Message = "duplicate username"
		return c.Status(fiber.ErrBadRequest.Code).JSON(data)
	}
	configs.Database.Create(&user)
	return c.Status(201).JSON(user)
}
