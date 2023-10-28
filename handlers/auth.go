package handlers

import (
	"context"
	"log"
	"os"

	"github.com/ricardonunosr/nait/data"
	"github.com/ricardonunosr/nait/db"

	"github.com/gofiber/fiber/v2"
	supa "github.com/nedpals/supabase-go"
)

func IsUserSignedIn(c *fiber.Ctx) bool {
	accessToken := c.Cookies("accessToken")
	if accessToken != "" {
		client := db.CreateSupabaseClient()
		if _, err := client.Auth.User(c.Context(), accessToken); err != nil {
			return false
		}
		return true
	}
	return false
}

func HandSignInView(c *fiber.Ctx) error {
	if IsUserSignedIn(c) {
		c.Redirect("/admin")
	}

	// TODO: for some reason the cookie is not being cleared properly
	c.ClearCookie("accessToken")
	return c.Render("signin", fiber.Map{})
}

func HandleSignIn(c *fiber.Ctx) error {
	admin := new(data.Admin)
	if err := c.BodyParser(admin); err != nil {
		return err
	}

	client := CreateSupabaseClient()
	resp, err := client.Auth.SignIn(context.Background(), supa.UserCredentials{
		Email:    admin.Email,
		Password: admin.Password,
	})

	if err != nil {
		log.Printf("Supabase auth error: %s\n", err)
		c.Set("Content-Type", "text/html")
		return c.SendString(`<p>Failed to login</p>`)
	}

	c.Cookie(&fiber.Cookie{
		Secure:   true,
		HTTPOnly: true,
		Name:     "accessToken",
		Value:    resp.AccessToken,
	})

	c.Append("HX-Redirect", "/admin")
	c.Set("Content-Type", "text/html")
	return c.SendString(`<p>Success to login</p>`)
}

func HandleSignOut(c *fiber.Ctx) error {
	client := CreateSupabaseClient()
	if err := client.Auth.SignOut(c.Context(), c.Cookies("accessToken")); err != nil {
		return err
	}

	c.ClearCookie("accessToken")
	return c.Redirect("/")
}

func CreateSupabaseClient() *supa.Client {
	return supa.CreateClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"))
}
