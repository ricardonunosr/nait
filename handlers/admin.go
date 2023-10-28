package handlers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/ricardonunosr/nait/data"
	"github.com/ricardonunosr/nait/db"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/checkout/session"
)

func HandleRegisterStaffView(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{})
}

func HandleRegisterStaff(c *fiber.Ctx) error {
	staff := new(data.Staff)

	if err := c.BodyParser(staff); err != nil {
		return err
	}

	if staff.Firstname != "" {
		db.DB.Exec(" INSERT INTO staff (username, email, firstname, lastname, password) VALUES(?,?,?,?,?)",
			staff.Username,
			staff.Email,
			staff.Firstname,
			staff.Lastname,
			staff.Password,
		)

		c.Set("Content-Type", "text/html")
		c.SendString(`<p>Registered Staff!</p>`)
		return c.Redirect("/admin")
	}

	c.Set("Content-Type", "text/html")
	// TODO: change this! this error messaging is awful for user
	return c.SendString(`<p>Something Went Wrong Registering Staff!</p>`)
}

func HandleStaffView(c *fiber.Ctx) error {
	if IsUserSignedIn(c) {
		rows, err := db.DB.Query("SELECT * FROM staff")

		if err != nil {
			log.Printf("Database error: %s\n", err)
			return c.Render("404", fiber.Map{})
		}

		var staffs []data.Staff

		for rows.Next() {
			var staff data.Staff
			rows.Scan(&staff.Username, &staff.Email, &staff.Firstname, &staff.Lastname, &staff.Password)
			staffs = append(staffs, staff)
		}

		stripe.Key = os.Getenv("STRIPE_KEY")

		params := &stripe.CheckoutSessionListParams{}
		params.Filters.AddFilter("limit", "", "3")
		i := session.List(params)

		promoters_details := &data.PromotersDetails{
			GuestCountSold:  make(map[string]int),
			GuestProfitSold: make(map[string]int),
		}
		for i.Next() {
			s := i.CheckoutSession()
			count_sold := promoters_details.GuestCountSold[s.ClientReferenceID]
			promoters_details.GuestCountSold[s.ClientReferenceID] = count_sold + 1

			profit_value := promoters_details.GuestProfitSold[s.ClientReferenceID]
			promoters_details.GuestProfitSold[s.ClientReferenceID] = int(s.AmountTotal) + profit_value
		}
		delete(promoters_details.GuestCountSold, "")
		delete(promoters_details.GuestProfitSold, "")

		var details []data.PromotersDetails2

		var detail data.PromotersDetails2
		for _, v := range promoters_details.GuestCountSold {
			detail.GuestCountSold = v
		}

		for _, v := range promoters_details.GuestProfitSold {
			detail.GuestProfitSold = v
		}
		details = append(details, detail)

		// Events
		rows, err = db.DB.Query("SELECT e.event_date, en.event_name FROM events e INNER JOIN events_name en ON e.event_id = en.event_id")

		if err != nil {
			return c.Render("404", fiber.Map{})
		}

		var events []data.Event

		for rows.Next() {
			var event data.Event
			rows.Scan(&event.EventDate, &event.EventName)
			events = append(events, event)
		}

		rows, err = db.DB.Query("SELECT * FROM events_name")
		if err != nil {
			return c.Render("404", fiber.Map{})
		}

		var event_names []data.EventName

		for rows.Next() {
			var event_name data.EventName
			rows.Scan(&event_name.EventID, &event_name.EventName)
			event_names = append(event_names, event_name)
		}

		return c.Render("admin", fiber.Map{
			"ClubName":         os.Getenv("CLUB_NAME"),
			"Staff":            staffs,
			"PromotersDetails": details,
			"Events":           events,
			"EventNames":       event_names,
		})
	}
	return c.Redirect("/admin/signin")
}
