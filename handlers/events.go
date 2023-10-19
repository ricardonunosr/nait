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
	event_date := c.FormValue("event_date")

	row, err := db.DB.Query("SELECT event_payment_url from events WHERE event_date = ?", event_date)

	if err != nil {
		return c.SendString("Something went wrong getting payment url...")
	}

	var event_payment_url string
	row.Next()
	row.Scan(&event_payment_url)

	c.Set("Content-Type", "text/html")
	return c.SendString(fmt.Sprintf(`<a href="%s">Buy</a>`, event_payment_url))
}

func HandleCreateNewEvent(c *fiber.Ctx) error {
	event_id := c.FormValue("event_name")
	event_date := c.FormValue("event_date")

	if event_id != "" {
		_, err := db.DB.Exec(" INSERT INTO events (event_id, event_date) VALUES(?,?)",
			event_id,
			event_date,
		)

		if err != nil {
			return c.SendString("Failed!")
		}

		row, err := db.DB.Query("SELECT event_name from events_name WHERE event_id = ?", event_id)
		if err != nil {
			return c.SendString("Something went wrong, could not find event name from event id...")
		}
		var event_name data.EventName
		row.Scan(&event_name.EventID, &event_name.EventName)

		stripe.Key = os.Getenv("STRIPE_KEY")
		product_params := &stripe.ProductParams{
			Name:        stripe.String(event_name.EventName + event_date),
			Description: stripe.String("10 euros"),
			Metadata:    map[string]string{"date": event_date},
		}
		starter_product, _ := product.New(product_params)

		price_params := &stripe.PriceParams{
			Currency:   stripe.String(string(stripe.CurrencyEUR)),
			Product:    stripe.String(starter_product.ID),
			UnitAmount: stripe.Int64(1000),
		}
		starter_price, _ := price.New(price_params)

		log.Println("Success! Here is your starter subscription product id: " + starter_product.ID)
		log.Println("Success! Here is your starter subscription price id: " + starter_price.ID)

		params := &stripe.PaymentLinkParams{
			Metadata: map[string]string{"date": event_date},
			AfterCompletion: &stripe.PaymentLinkAfterCompletionParams{
				Type: stripe.String("redirect"),
				Redirect: &stripe.PaymentLinkAfterCompletionRedirectParams{
					URL: stripe.String(fmt.Sprintf("http://localhost:8080/guest/completed?session_id={CHECKOUT_SESSION_ID}&event_date=%s", event_date)),
				},
			},
			LineItems: []*stripe.PaymentLinkLineItemParams{
				{
					Price:    stripe.String(starter_price.ID),
					Quantity: stripe.Int64(1),
				},
			},
		}
		pl, _ := paymentlink.New(params)

		_, err = db.DB.Exec("UPDATE events SET event_payment_url = ? WHERE event_date = ?",
			pl.URL,
			event_date,
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
	event_name := new(data.EventName)

	if err := c.BodyParser(event_name); err != nil {
		return err
	}

	if event_name.EventName != "" {
		db.DB.Exec(" INSERT INTO events_name (event_name) VALUES(?)",
			event_name.EventName,
		)

		c.Set("Content-Type", "text/html")
		return c.SendString(`<p>Registered Event Name!</p>`)
	}

	c.Set("Content-Type", "text/html")
	// TODO: change this! this error messaging is awful for user
	return c.SendString(`<p>Something Went Wrong Registering Event name!</p>`)
}
