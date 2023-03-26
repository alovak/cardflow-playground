package acquirer

type Card struct {
	Number         string
	ExpirationDate string
	CVV            string
}

type SafeCard struct {
	First6         string
	Last4          string
	ExpirationDate string
}
