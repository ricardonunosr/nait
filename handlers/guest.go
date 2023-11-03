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
	guestName := c.Params("username")
	rows, err := db.Db.Query("SELECT * FROM staff WHERE username = ?", guestName)

	if err != nil {
		log.Println("[ERROR] Did not find ", guestName, " ...")
		return c.Render("404", fiber.Map{})
	}

	var staff data.Staff
	for rows.Next() {
		rows.Scan(&staff.Username, &staff.Email, &staff.Firstname, &staff.Lastname, &staff.Password)
	}

	if staff.Username == "" {
		return c.Render("404", fiber.Map{})
	}

	rows, err = db.Db.Query("SELECT en.event_name, e.event_date, e.event_payment_url FROM events e INNER JOIN events_name en ON e.event_id = en.event_id LIMIT 10")

	if err != nil {
		return c.SendString("Could not find events!")
	}

	next5events := make([]data.Event, 5, 5)
	for rows.Next() {
		var event data.Event
		rows.Scan(&event.EventName, &event.EventDate, &event.EventURL)
		convDate, err := time.Parse(time.RFC3339, event.EventDate)
		if err != nil {
			c.SendString("Failed conversion of date")
		}
		event.EventDate = convDate.Format("2006-01-02")
		next5events = append(next5events, event)
	}

	// Get the stripe url for the first select
	// TODO: handle this better!
	if len(next5events) != 0 {
		firstEventUrl := next5events[0].EventURL
		_ = firstEventUrl
	}

	return c.Render("pay", fiber.Map{
		"GuestUsername":  staff.Username,
		"GuestFirstname": staff.Firstname,
		"GuestLastname":  staff.Lastname,
		"Next5Events":    next5events,
		// "EventURL":       first_event_url,
	})
}

func HandleCompletedView(c *fiber.Ctx) error {
	stripe.Key = os.Getenv("STRIPE_KEY")
	sessionId := c.Query("session_id")
	eventDate := c.Query("event_date")
	s, _ := session.Get(sessionId, nil)

	// Check if there is a `code` already generated for the sessionId and eventDate
	var code data.Code
	db.Db.QueryRow("SELECT * FROM codes WHERE checkout_session_id = ? AND event_date = ?", sessionId, eventDate).Scan(&code.EventDate, &code.Code, &code.CheckoutSessionId)
	if code.Code != "" {
		log.Println("Already a code generated for this checkout session and event date...")
		return c.Render("404", fiber.Map{
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

		_, err = db.Db.Exec("INSERT INTO codes (code, event_date, checkout_session_id) VALUES(?,?,?)", otp, eventDate, sessionId)
		if err != nil {
			log.Println("Something went wrong when inserting OTP...")
			return c.Render("404", fiber.Map{})
		}

		return c.Render("completed", fiber.Map{})
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
