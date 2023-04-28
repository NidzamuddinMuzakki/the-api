package controllers

import (
	"log"
	"strconv"

	"github.com/NidzamuddinMuzakki/the-api/configs"
	"github.com/NidzamuddinMuzakki/the-api/models"
	"github.com/NidzamuddinMuzakki/the-api/responses"
	"github.com/gofiber/fiber/v2"
)

func GetArticles(c *fiber.Ctx) error {
	var total int64
	perpage, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("offset"))
	var articles []models.Article
	log.Println(page, perpage)
	configs.Database.Find(&articles).Count(&total)
	configs.Database.Limit(perpage).Offset((page - 1) * perpage).Order("id ASC").Find(&articles)
	data := responses.ArticleResponse{
		Status:    200,
		Message:   "success",
		Data:      articles,
		TotalData: total,
	}
	return c.Status(200).JSON(data)
}
func GetArticle(c *fiber.Ctx) error {
	var article models.Article
	id := c.Params("id")
	result := configs.Database.Find(&article, "id=?", id)
	data := responses.UserResponse{
		Status:  404,
		Message: "not found",
		Data:    nil,
	}
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(data)
	}
	data.Status = 200
	data.Message = "success"
	data.Data = article

	return c.Status(200).JSON(data)
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
