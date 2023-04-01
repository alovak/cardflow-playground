package authorizer

type Authorizer interface {
	AuthorizeRequest(req AuthorizationRequest) (AuthorizationResponse, error)
}

type AuthorizationRequest struct {
	Amount   int
	Currency string
	Card     Card
}

type AuthorizationResponse struct {
	AuthorizationCode string
	ApprovalCode      string
}

type Card struct {
	ID                    string
	AccountID             string
	Number                string
	ExpirationDate        string
	CardVerificationValue string
}
