package handlers

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go-auth-api/models"
	"go-auth-api/utils"
)

func getUserCollection() *mongo.Collection {
	return utils.DB.Collection("users")
}

func getUserFromToken(c *fiber.Ctx) (*models.User, error) {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Missing or invalid token.")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token.")
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token.")
	}

	var user models.User
	collection := getUserCollection()
	err = collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found.")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "An error occurred while fetching user data.")
	}

	return &user, nil
}

func GetUser(c *fiber.Ctx) error {
	user, err := getUserFromToken(c)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input."})
	}

	collection := getUserCollection()

	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username already exists."})
	} else if err != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An error occurred while processing your request."})
	}

	user.Role = "user"
	user.Provider = "credentials"
	user.ID = uuid.NewString()
	user.Password = utils.HashPassword(user.Password)
	user.ProfileURL = "https://www.gravatar.com/avatar/" + user.ID

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An error occurred while processing your request."})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully."})
}

func Login(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input."})
	}

	var user models.User
	collection := getUserCollection()
	err := collection.FindOne(context.Background(), bson.M{"username": input.Username}).Decode(&user)
	if err != nil || !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials."})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Logged out successfully."})
}

func Admin(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Welcome, Admin!"})
}
