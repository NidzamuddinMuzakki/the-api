package controllers

import (
	"errors"
	"fmt"
	"log"

	"github.com/NidzamuddinMuzakki/the-api/configs"
	"github.com/NidzamuddinMuzakki/the-api/models"
	"github.com/NidzamuddinMuzakki/the-api/requests"
	"github.com/NidzamuddinMuzakki/the-api/utils"

	"github.com/NidzamuddinMuzakki/the-api/responses"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetArticles(c *fiber.Ctx) error {
	var total int64
	pagination := new(requests.Pagination)
	data := responses.ErrResponse{
		Status:  200,
		Message: "success",
	}
	if err := c.QueryParser(pagination); err != nil {
		data.Status = fiber.ErrBadRequest.Code
		data.Message = "err bad request"
		data.Error = err.Error()
		return c.Status(fiber.ErrBadRequest.Code).JSON(data)
	}
	validate := validator.New()
	err := validate.Struct(pagination)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]responses.ApiError, len(ve))
			for i, fe := range ve {
				out[i] = responses.ApiError{Field: fe.Field(), Msg: utils.MsgForTag(fe.Tag())}
			}
			data.Status = fiber.ErrBadRequest.Code
			data.Message = "bad request"
			data.Error = out
		}
		return c.Status(fiber.ErrBadRequest.Code).JSON(data)
	}
	// fmt.Println(err.Error())

	var articles []models.Article
	log.Println(pagination.Limit, pagination.Offset)
	configs.Database.Find(&articles).Count(&total)
	configs.Database.Limit(pagination.Limit).Offset((pagination.Offset - 1) * pagination.Limit).Order("id ASC").Find(&articles)
	dataArticle := responses.ArticleResponse{
		Status:    200,
		Message:   "success",
		TotalData: total,
		Data:      articles,
	}
	return c.Status(200).JSON(dataArticle)
}
func GetArticle(c *fiber.Ctx) error {
	var article models.Article
	req := new(requests.RequestAData)
	data := responses.ErrResponse{
		Status:  200,
		Message: "success",
	}
	err := c.ParamsParser(req)
	if err != nil {
		data.Status = fiber.ErrBadRequest.Code
		data.Message = "err bad request"
		data.Error = err.Error()
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		fmt.Println(err.Error())
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {

			out := make([]responses.ApiError, len(ve))
			for i, fe := range ve {
				out[i] = responses.ApiError{Field: fe.Field(), Msg: utils.MsgForTag(fe.Tag())}
			}
			data.Status = fiber.ErrBadRequest.Code
			data.Message = "bad request"
			data.Error = out
		}
		return c.Status(fiber.ErrBadRequest.Code).JSON(data)
	}
	result := configs.Database.Find(&article, "id=?", req.Id)
	data.Status = 404
	data.Message = "not found"
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(data)
	}
	dataArticle := responses.UserResponse{
		Status:  200,
		Message: "success",
		Data:    article,
	}

	return c.Status(200).JSON(dataArticle)
}

func AddArticle(c *fiber.Ctx) error {
	user := new(models.Article)
	data := responses.UserResponse{
		Status: fiber.ErrBadRequest.Code,

		Data: nil,
	}
	if err := c.BodyParser(user); err != nil {
		data.Message = err.Error()
		return c.Status(fiber.ErrBadRequest.Code).JSON(data)
	}

	configs.Database.Create(&user)
	return c.Status(201).JSON(user)
}
