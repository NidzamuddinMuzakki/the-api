package middlewares

import (
	"context"
	"fmt"
	"strings"

	"github.com/NidzamuddinMuzakki/the-api/configs"
	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2"
)

func DeserializeUser(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")
	url := c.Request().URI()
	if strings.Contains(url.String(), "/token/refresh") || strings.Contains(url.String(), "/token") || strings.Contains(url.String(), "/register") {
		return c.Next()
	}
	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	tokenByte, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])

		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("invalidate token: %v", err)})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})

	}
	ss, err := configs.RedisClient.Get(context.TODO(), claims["uuid"].(string)).Result()
	if err != redis.Nil && ss != "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "your token in black list data"})
	}
	// var user models.User
	// initializers.DB.First(&user, "id = ?", fmt.Sprint(claims["sub"]))

	// if user.ID.String() != claims["sub"] {
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	// }

	// c.Locals("user", models.FilterUserRecord(&user))

	return c.Next()
}
