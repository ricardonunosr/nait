package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ricardonunosr/nait/data"
	"github.com/ricardonunosr/nait/db"
)

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

	rows, err = db.DB.Query("SELECT en.event_name, e.event_date, e.event_payment_url FROM events e INNER JOIN events_name en ON e.event_id = en.event_id LIMIT 5")

	if err != nil {
		return c.SendString("Could not find events!")
	}

	// TODO: Could make static
	var next_5_events []data.Event
	for rows.Next() {
		var event data.Event
		rows.Scan(&event.EventName, &event.EventDate, &event.EventURL)
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
