package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"go-auth-api/models"
)

var googleOauthConfig *oauth2.Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if redirectURL == "" || clientID == "" || clientSecret == "" {
		log.Fatal("Required environment variables are not set")
	}

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

func OAuthGoogleLogin(c *fiber.Ctx) error {
	url := googleOauthConfig.AuthCodeURL("random_state")
	return c.Redirect(url)
}

func OAuthGoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to get token")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to get user info")
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to decode user info")
	}

	email := userInfo["email"].(string)
	googleID := userInfo["id"].(string)

	var user models.User
	collection := getUserCollection()
	err = collection.FindOne(context.Background(), bson.M{"email": email, "provider": "google"}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		user = models.User{
			ID:         uuid.NewString(),
			Email:      email,
			Provider:   "google",
			ProviderID: googleID,
			Role:       "user",
			Username:   email,
			ProfileURL: userInfo["picture"].(string),
		}
		_, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create user.")
		}
	}

	tokenJWT := jwt.New(jwt.SigningMethodHS256)
	claims := tokenJWT.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := tokenJWT.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}
