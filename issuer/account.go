package issuer

type CreateAccount struct {
	Balance  int
	Currency string
}

type Account struct {
	ID       string
	Balance  int
	Currency string
}
