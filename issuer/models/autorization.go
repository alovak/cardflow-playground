package models

type AuthorizationRequest struct {
	Amount   int
	Currency string
	Card     Card
	Merchant Merchant
}

type AuthorizationResponse struct {
	AuthorizationCode string
	ApprovalCode      string
}
