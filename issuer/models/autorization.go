package models

type AuthorizationRequest struct {
	Amount   int64
	Currency string
	Card     Card
	Merchant Merchant
}

type AuthorizationResponse struct {
	AuthorizationCode string
	ApprovalCode      string
}
