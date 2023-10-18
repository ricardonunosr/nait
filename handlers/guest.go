package handlers

import (
	"crypto/rand"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ricardonunosr/nait/data"
	"github.com/ricardonunosr/nait/db"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/checkout/session"
)

const otpChars = "1234567890"

func HandleGuestView(c *fiber.Ctx) error {
	guest_name := c.Params("username")
	rows, err := db.DB.Query("SELECT * FROM staff WHERE username = ?", guest_name)

	if err != nil {
		log.Println("[ERROR] Did not find ", guest_name, " ...")
		return c.Render("404", fiber.Map{})
	}

	var staff data.Staff
	for rows.Next() {
		rows.Scan(&staff.Username, &staff.Email, &staff.Firstname, &staff.Lastname, &staff.Password)
	}

	if staff.Username == "" {
		return c.Render("404", fiber.Map{})
	}

	rows, err = db.DB.Query("SELECT en.event_name, e.event_date, e.event_payment_url FROM events e INNER JOIN events_name en ON e.event_id = en.event_id LIMIT 10")

	if err != nil {
		return c.SendString("Could not find events!")
	}

	// TODO: Could make static
	var next_5_events []data.Event
	for rows.Next() {
		var event data.Event
		rows.Scan(&event.EventName, &event.EventDate, &event.EventURL)
		conv_date, err := time.Parse(time.RFC3339, event.EventDate)
		if err != nil {
			c.SendString("Failed conversion of date")
		}
		event.EventDate = conv_date.Format("2006-01-02")
		next_5_events = append(next_5_events, event)
	}

	// Get the stripe url for the first select
	first_event_url := next_5_events[0].EventURL

	return c.Render("pay", fiber.Map{
		"GuestUsername":  staff.Username,
		"GuestFirstname": staff.Firstname,
		"GuestLastname":  staff.Lastname,
		"Next5Events":    next_5_events,
		"EventURL":       first_event_url,
	})
}

func HandleCompletedView(c *fiber.Ctx) error {
	stripe.Key = os.Getenv("STRIPE_KEY")
	session_id := c.Query("session_id")
	event_date := c.Query("event_date")
	s, _ := session.Get(session_id, nil)

	// Check if there is a `code` already generated for the session_id and event_date
	var code data.Code
	db.DB.QueryRow("SELECT * FROM codes WHERE checkout_session_id = ? AND event_date = ?", session_id, event_date).Scan(&code.EventDate, &code.Code, &code.CheckoutSessionId)
	if code.Code != "" {
		log.Println("Already a code generated for this checkout session and event date...")
		return c.Render("completed", fiber.Map{
			"ClubName": os.Getenv("CLUB_NAME"),
		})
	}

	if s.Status == "complete" {
		log.Println("Sent Email...")
		otp, err := GenerateOTP(6)
		if err != nil {
			log.Println("Something went wrong when generating OTP...")
			return c.Render("404", fiber.Map{})
		}

		_, err = db.DB.Exec("INSERT INTO codes (code, event_date, checkout_session_id) VALUES(?,?,?)", otp, event_date, session_id)
		if err != nil {
			log.Println("Something went wrong when inserting OTP...")
			return c.Render("404", fiber.Map{})
		}

		return c.Render("completed", fiber.Map{
			"ClubName": os.Getenv("CLUB_NAME"),
		})
	} else {
		log.Println("Did not find a checkout session id...")
		return c.Render("404", fiber.Map{})
	}

}

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
