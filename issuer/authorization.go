package issuer

type AuthorizationRequest struct {
	Amount   int
	Currency string
	Card     Card
}

type AuthorizationResponse struct {
	AuthorizationCode string
	ApprovalCode      string
}
