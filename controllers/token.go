package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NidzamuddinMuzakki/the-api/configs"
	"github.com/NidzamuddinMuzakki/the-api/requests"
	"github.com/NidzamuddinMuzakki/the-api/responses"
	"github.com/NidzamuddinMuzakki/the-api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2"
)

type Token struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshtoken"`
}

func GetToken(c *fiber.Ctx) error {
	var user requests.LoginRequest
	dataErr := responses.ErrResponse{
		Status:  200,
		Message: "success",
	}
	if err := c.BodyParser(&user); err != nil {
		dataErr.Status = fiber.ErrBadRequest.Code
		dataErr.Message = "bad request"
		dataErr.Error = err.Error()
		return c.Status(fiber.ErrBadRequest.Code).JSON(dataErr)
	}
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		fmt.Println(err.Error())
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {

			out := make([]responses.ApiError, len(ve))
			for i, fe := range ve {
				out[i] = responses.ApiError{Field: fe.Field(), Msg: utils.MsgForTag(fe.Tag())}
			}
			dataErr.Status = fiber.ErrBadRequest.Code
			dataErr.Message = "bad request"
			dataErr.Error = out
		}
		return c.Status(fiber.ErrBadRequest.Code).JSON(dataErr)
	}
	result := configs.Database.Table("users").Find(&user, "username=? and password=?", user.Username, user.Password)

	if result.RowsAffected == 0 {
		dataErr.Status = 401
		dataErr.Message = "Unauthorized"
		return c.Status(401).JSON(dataErr)
	}

	tokens, err := configs.GenerateTokenPair(user.Username)
	if err != nil {
		dataErr.Status = 401
		dataErr.Message = "Unauthorized"
		dataErr.Error = err.Error()
		return c.Status(401).JSON(dataErr)
	}
	token := Token{
		Token:        tokens["token"],
		RefreshToken: tokens["refreshtoken"],
	}
	data := responses.UserResponse{
		Status:  200,
		Message: "success",
		Data:    token,
	}

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
		return c.Status(401).JSON(err.Error() + "unauthorized")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if int(claims["sub"].(float64)) == 1 && claims["name"] != nil && claims["uuid"] != nil {
			ss, errS := configs.RedisClient.Get(context.TODO(), claims["uuid"].(string)).Result()
			// fmt.Println(ss, errS, ss == "", ss != "", errS != nil, !(ss == ""))
			if errS != redis.Nil || ss != "" {
				return c.Status(401).JSON("refresh token in black list data")
			}

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
			str, _ := configs.RedisClient.Set(context.TODO(), claims["uuid"].(string), claims["uuid"].(string), time.Hour*24).Result()

			if str != "OK" {
				return c.Status(401).JSON("unauthorized")
			}
			return c.Status(200).JSON(data)
		}

		return c.Status(401).JSON("unauthorized")
	}
	return c.Status(401).JSON("unauthorized")
}
