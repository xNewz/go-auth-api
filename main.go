package main

import (
	"go-auth-api/handlers"
	"go-auth-api/middlewares"
	"go-auth-api/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	utils.ConnectDB(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_DB"))

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
		AllowMethods: "GET,POST,PUT",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)
	app.Post("/logout", handlers.Logout)

	app.Get("/oauth/google", handlers.OAuthGoogleLogin)
	app.Get("/oauth/google/callback", handlers.OAuthGoogleCallback)

	app.Get("/user", middlewares.AuthMiddleware, handlers.GetUser)
	app.Get("/admin", middlewares.AdminMiddleware, handlers.Admin)
	app.Put("/admin/edit-role/:user_id", middlewares.AdminMiddleware, handlers.EditUserRole)

	log.Fatal(app.Listen(":3000"))
}
