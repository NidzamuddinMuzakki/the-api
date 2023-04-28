package main

import (
	"log"

	"github.com/NidzamuddinMuzakki/the-api/configs"
	"github.com/NidzamuddinMuzakki/the-api/controllers"
	"github.com/NidzamuddinMuzakki/the-api/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	configs.Connect()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(middlewares.DeserializeUser)
	app.Post("/api/api/token", controllers.GetToken)
	app.Post("/api/api/token/refresh", controllers.RefreshToken)
	app.Post("/api/register", controllers.AddUser)
	app.Post("/api/article", controllers.AddArticle)
	app.Get("/api/article", controllers.GetArticles)
	app.Get("/api/article/:id", controllers.GetArticle)

	app.Get("/api/profile", controllers.GetUsers)
	app.Get("/uploads/:name", func(c *fiber.Ctx) error {
		c.SendFile("uploads/" + c.Params("name"))
		return nil
	})
	log.Fatal(app.Listen(":8080"))
}
