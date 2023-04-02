package models

type Card struct {
	Number                string
	ExpirationDate        string
	CardVerificationValue string
}

type SafeCard struct {
	First6         string
	Last4          string
	ExpirationDate string
}
