package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ricardonunosr/nait/db"
	"github.com/ricardonunosr/nait/handlers"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/django/v3"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()

	engine := django.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(favicon.New(favicon.Config{
		File: "./static/favicon.ico",
		URL:  "/favicon.ico",
	}))
	app.Static("/", "./static")

	// Main Page ---------------------------
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"ClubName": os.Getenv("CLUB_NAME"),
		})
	})

	// Admin ------------------------------
	admin := app.Group("/admin")
	admin.Get("/", handlers.HandleStaffView)
	admin.Get("/signin", handlers.HandSignInView)
	admin.Post("/signin", handlers.HandleSignIn)
	admin.Get("/staff/new", handlers.HandleRegisterStaffView)
	admin.Post("/staff", handlers.HandleRegisterStaff)

	// Guest ------------------------------
	guest := app.Group("/guest")
	guest.Get("/completed", handlers.HandleCompletedView)
	guest.Get("/:username", handlers.HandleGuestView)

	// Event ------------------------------
	event := app.Group("/event")
	event.Get("/", handlers.HandleGetEvent)
	event.Post("/check", handlers.HandleCheckCode)
	event.Post("/new", handlers.HandleCreateNewEvent)
	event.Post("/name/new", handlers.HandleCreateNewEventName)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
