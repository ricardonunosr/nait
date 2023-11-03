package handlers

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/ricardonunosr/nait/data"
	"github.com/ricardonunosr/nait/db"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/paymentlink"
	"github.com/stripe/stripe-go/v75/price"
	"github.com/stripe/stripe-go/v75/product"
)

func HandleGetEvent(c *fiber.Ctx) error {
	eventDate := c.FormValue("event_date")

	row, err := db.Db.Query("SELECT event_payment_url from events WHERE event_date = ?", eventDate)

	if err != nil {
		return c.SendString("Something went wrong getting payment url...")
	}

	var eventPaymentUrl string
	row.Next()
	row.Scan(&eventPaymentUrl)

	c.Set("Content-Type", "text/html")
	return c.SendString(fmt.Sprintf(`<a href="%s">Buy</a>`, eventPaymentUrl))
}

func HandleCreateNewEvent(c *fiber.Ctx) error {
	eventId := c.FormValue("event_name")
	eventDate := c.FormValue("event_date")

	if eventId != "" {
		_, err := db.Db.Exec(" INSERT INTO events (event_id, event_date) VALUES(?,?)",
			eventId,
			eventDate,
		)

		if err != nil {
			return c.SendString("Failed!")
		}

		row, err := db.Db.Query("SELECT event_name from events_name WHERE event_id = ?", eventId)
		if err != nil {
			return c.SendString("Something went wrong, could not find event name from event id...")
		}
		var eventName data.EventName
		row.Scan(&eventName.EventID, &eventName.EventName)

		stripe.Key = os.Getenv("STRIPE_KEY")
		productParams := &stripe.ProductParams{
			Name:        stripe.String(eventName.EventName + eventDate),
			Description: stripe.String("10 euros"),
			Metadata:    map[string]string{"date": eventDate},
		}
		starterProduct, _ := product.New(productParams)

		priceParams := &stripe.PriceParams{
			Currency:   stripe.String(string(stripe.CurrencyEUR)),
			Product:    stripe.String(starterProduct.ID),
			UnitAmount: stripe.Int64(1000),
		}
		starterPrice, _ := price.New(priceParams)

		log.Println("Success! Here is your starter subscription product id: " + starterProduct.ID)
		log.Println("Success! Here is your starter subscription price id: " + starterPrice.ID)

		params := &stripe.PaymentLinkParams{
			Metadata: map[string]string{"date": eventDate},
			AfterCompletion: &stripe.PaymentLinkAfterCompletionParams{
				Type: stripe.String("redirect"),
				Redirect: &stripe.PaymentLinkAfterCompletionRedirectParams{
					URL: stripe.String(fmt.Sprintf("http://localhost:8080/guest/completed?session_id={CHECKOUT_SESSION_ID}&event_date=%s", eventDate)),
				},
			},
			LineItems: []*stripe.PaymentLinkLineItemParams{
				{
					Price:    stripe.String(starterPrice.ID),
					Quantity: stripe.Int64(1),
				},
			},
		}
		pl, _ := paymentlink.New(params)

		_, err = db.Db.Exec("UPDATE events SET event_payment_url = ? WHERE event_date = ?",
			pl.URL,
			eventDate,
		)

		if err != nil {
			c.SendString("Something went wrong updating url for the event!")
		}

		log.Println(pl.URL)

		return c.SendString("Success")
	}

	return c.SendString("Failed")
}

func HandleCreateNewEventName(c *fiber.Ctx) error {
	eventName := new(data.EventName)

	if err := c.BodyParser(eventName); err != nil {
		return err
	}

	if eventName.EventName != "" {
		_, err := db.Db.Exec(" INSERT INTO events_name (event_name) VALUES(?)",
			eventName.EventName,
		)

		if err != nil {
			log.Printf("Database error: %s", err)
			return c.SendString(`Something went wrong!`)
		}

		c.Set("Content-Type", "text/html")
		return c.SendString(`<p>Registered Event Name!</p>`)
	}

	c.Set("Content-Type", "text/html")
	// TODO: change this! this error messaging is awful for user
	return c.SendString(`<p>Something Went Wrong Registering Event name!</p>`)
}

func HandleCheckCode(c *fiber.Ctx) error {
	guestCode := c.FormValue("guest_code")
	eventDate := "2023-11-01"

	if guestCode != "" {
		var eventCode data.Code
		db.Db.QueryRow("SELECT * FROM codes WHERE code = ? AND event_date = ?", guestCode, eventDate).Scan(&eventCode.EventDate, &eventCode.Code, &eventCode.CheckoutSessionId)
		if eventCode.Code != "" {
			return c.SendString("Checked code successfully")
		} else {
			return c.SendString("Failed to check code...")
		}
	}

	return c.SendString("Failed getting code...")
}
