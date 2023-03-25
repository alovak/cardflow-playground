package issuer

type CreateAccountRequest struct {
	Balance  int
	Currency string
}

type Account struct {
	ID       string
	Balance  int
	Currency string
}
