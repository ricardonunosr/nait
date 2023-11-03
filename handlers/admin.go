package handlers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/ricardonunosr/nait/data"
	"github.com/ricardonunosr/nait/db"
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
		_, err := db.Db.Exec(" INSERT INTO staff (username, email, firstname, lastname, password) VALUES(?,?,?,?,?)",
			staff.Username,
			staff.Email,
			staff.Firstname,
			staff.Lastname,
			staff.Password,
		)

		if err != nil {
			log.Printf("Database error: %s", err)
			return c.SendString("Something went wrong")
		}

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
		rows, err := db.StaffList.Query()
		// rows, err := db.DB.CONN.Query("SELECT * FROM staff")

		if err != nil {
			log.Printf("Database error: %s\n", err)
			return c.Render("404", fiber.Map{})
		}

		// var staffs []data.Staff
		// NOTE: MAX 10 Staff
		staffs := make([]data.Staff, 0, 10)

		for rows.Next() {
			var staff data.Staff
			rows.Scan(&staff.Username, &staff.Email, &staff.Firstname, &staff.Lastname, &staff.Password)
			staffs = append(staffs, staff)
		}

		// stripe.Key = os.Getenv("STRIPE_KEY")

		// params := &stripe.CheckoutSessionListParams{}
		// params.Filters.AddFilter("limit", "", "3")
		// i := session.List(params)

		// promoters_details := &data.PromotersDetails{
		// 	GuestCountSold:  make(map[string]int),
		// 	GuestProfitSold: make(map[string]int),
		// }
		// for i.Next() {
		// 	s := i.CheckoutSession()
		// 	count_sold := promoters_details.GuestCountSold[s.ClientReferenceID]
		// 	promoters_details.GuestCountSold[s.ClientReferenceID] = count_sold + 1

		// 	profit_value := promoters_details.GuestProfitSold[s.ClientReferenceID]
		// 	promoters_details.GuestProfitSold[s.ClientReferenceID] = int(s.AmountTotal) + profit_value
		// }
		// delete(promoters_details.GuestCountSold, "")
		// delete(promoters_details.GuestProfitSold, "")

		// var details []data.PromotersDetails2

		// var detail data.PromotersDetails2
		// for _, v := range promoters_details.GuestCountSold {
		// 	detail.GuestCountSold = v
		// }

		// for _, v := range promoters_details.GuestProfitSold {
		// 	detail.GuestProfitSold = v
		// }
		// details = append(details, detail)

		// Events
		rows, err = db.EventsList.Query()

		if err != nil {
			return c.Render("404", fiber.Map{})
		}

		// var events []data.Event
		events := make([]data.Event, 0, 10)

		for rows.Next() {
			var event data.Event
			rows.Scan(&event.EventDate, &event.EventName)
			events = append(events, event)
		}

		rows, err = db.EventsNameList.Query()
		if err != nil {
			return c.Render("404", fiber.Map{})
		}

		// var eventNames []data.EventName
		eventNames := make([]data.EventName, 0, 10)

		for rows.Next() {
			var eventName data.EventName
			rows.Scan(&eventName.EventID, &eventName.EventName)
			eventNames = append(eventNames, eventName)
		}

		return c.Render("admin", fiber.Map{
			"ClubName": os.Getenv("CLUB_NAME"),
			"Staff":    staffs,
			// "PromotersDetails": details,
			"Events":     events,
			"EventNames": eventNames,
		})
	}
	return c.Redirect("/admin/signin")
}
