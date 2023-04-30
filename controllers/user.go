package controllers

import (
	"errors"
	"fmt"

	"github.com/NidzamuddinMuzakki/the-api/configs"
	"github.com/NidzamuddinMuzakki/the-api/models"
	"github.com/NidzamuddinMuzakki/the-api/responses"
	"github.com/NidzamuddinMuzakki/the-api/utils"
	"github.com/go-playground/validator/v10"
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
	dataErr := responses.ErrResponse{
		Status:  fiber.ErrBadRequest.Code,
		Message: "bad request",
		Error:   nil,
	}
	if err := c.BodyParser(user); err != nil {

		dataErr.Error = err.Error()
		return c.Status(fiber.ErrBadRequest.Code).JSON(dataErr)
	}
	validate := validator.New()
	validate.RegisterValidation("regexp", utils.Regexp)
	err := validate.Struct(user)
	if err != nil {
		fmt.Println(err.Error())
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {

			out := make([]responses.ApiError, len(ve))
			for i, fe := range ve {
				out[i] = responses.ApiError{Field: fe.Field(), Msg: utils.MsgForTag(fe.Tag())}
			}

			dataErr.Error = out
		}
		return c.Status(fiber.ErrBadRequest.Code).JSON(dataErr)
	}
	file, err := c.FormFile("profile_image")
	if err != nil {
		dataErr.Error = err.Error()
		return c.Status(fiber.ErrBadRequest.Code).JSON(dataErr)

	}
	user.ProfileImage = "uploads/" + id.String() + ".jpg"
	c.SaveFile(file, user.ProfileImage)
	// user.ProfileImage = "uploads/" + id.String() + ".jpg"

	// fmt.Println(user.ProfileImage, user.FirstName, user.Las)
	datas := configs.Database.First(&user, "username=?", user.Username)
	if datas.RowsAffected != 0 {
		dataErr.Message = "duplicate username"
		return c.Status(fiber.ErrBadRequest.Code).JSON(dataErr)
	}
	configs.Database.Create(&user)
	return c.Status(201).JSON(user)
}
