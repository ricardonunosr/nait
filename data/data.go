package data

type Admin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type Staff struct {
	Username  string
	Email     string
	Firstname string
	Lastname  string
	Password  string
}

type PromotersDetails struct {
	GuestCountSold  map[string]int
	GuestProfitSold map[string]int
}

type PromotersDetails2 struct {
	GuestCountSold  int
	GuestProfitSold int
}

type EventName struct {
	EventID   int
	EventName string
}

type Event struct {
	EventName string
	EventDate string
	EventURL  string
}

type Code struct {
	EventDate         string
	Code              string
	CheckoutSessionId string
}
