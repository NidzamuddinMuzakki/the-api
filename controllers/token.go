package controllers

import (
	"fmt"

	"github.com/NidzamuddinMuzakki/the-api/configs"
	"github.com/NidzamuddinMuzakki/the-api/models"
	"github.com/NidzamuddinMuzakki/the-api/responses"
	"github.com/dgrijalva/jwt-go"

	"github.com/gofiber/fiber/v2"
)

type Token struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshtoken"`
}

func GetToken(c *fiber.Ctx) error {
	var user models.User
	data := responses.UserResponse{
		Status: fiber.ErrBadRequest.Code,

		Data: nil,
	}
	if err := c.BodyParser(&user); err != nil {
		data.Message = err.Error()
		return c.Status(fiber.ErrBadRequest.Code).JSON(data)
	}
	result := configs.Database.Find(&user, "username=? and password=?", user.Username, user.Password)

	if result.RowsAffected == 0 {
		data.Status = 401
		data.Message = "Unauthorized"
		return c.Status(401).JSON(data)
	}
	tokens, err := configs.GenerateTokenPair(user.Username)
	if err != nil {
		data.Status = 401
		data.Message = "Unauthorized"
		return c.Status(401).JSON(data)
	}
	data.Status = 200
	data.Message = "success"
	token := Token{
		Token:        tokens["token"],
		RefreshToken: tokens["refreshtoken"],
	}
	data.Data = token

	return c.Status(200).JSON(data)
}

func RefreshToken(c *fiber.Ctx) error {
	var tokens Token
	data := responses.UserResponse{
		Status: fiber.ErrBadRequest.Code,

		Data: nil,
	}
	if err := c.BodyParser(&tokens); err != nil {
		data.Message = err.Error()
		return c.Status(fiber.ErrBadRequest.Code).JSON(data)
	}
	token, err := jwt.Parse(tokens.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])

		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if err != nil {
		return c.Status(401).JSON("unauthorized")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if int(claims["sub"].(float64)) == 1 && claims["name"] != nil {

			newTokenPair, err := configs.GenerateTokenPair(claims["name"].(string))
			if err != nil {
				return c.Status(401).JSON("unauthorized")
			}
			data := responses.UserResponse{}
			data.Status = 200
			data.Message = "success"
			token := Token{
				Token:        newTokenPair["token"],
				RefreshToken: newTokenPair["refreshtoken"],
			}
			data.Data = token
			return c.Status(200).JSON(data)
		}

		return c.Status(401).JSON("unauthorized")
	}
	return c.Status(401).JSON("unauthorized")
}
