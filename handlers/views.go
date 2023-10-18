package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func HandleLogInView(c *fiber.Ctx) error {
	if IsUserSignedIn(c) {
		c.Redirect("/admin")
	}

	// TODO: for some reason the cookie is not being cleared properly
	c.ClearCookie("accessToken")
	return c.Render("login", fiber.Map{
		"ClubName": os.Getenv("CLUB_NAME"),
	})
}
