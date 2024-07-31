package middlewares

import (
	"context"
	"os"

	"go-auth-api/models"
	"go-auth-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserCollection() *mongo.Collection {
	return utils.DB.Collection("users")
}

func parseToken(c *fiber.Ctx) (*jwt.Token, jwt.MapClaims, error) {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "Missing or invalid token.")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token.")
	}

	return token, claims, nil
}

func AuthMiddleware(c *fiber.Ctx) error {
	_, claims, err := parseToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	c.Locals("user_id", claims["id"])
	c.Locals("user_role", claims["role"])

	return c.Next()
}

func AdminMiddleware(c *fiber.Ctx) error {
	_, claims, err := parseToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token."})
	}

	var user models.User
	collection := getUserCollection()
	err = collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found."})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied."})
	}

	c.Locals("role", user.Role)

	return c.Next()
}
