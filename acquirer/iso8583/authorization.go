package iso8583

import "github.com/moov-io/iso8583/field"

type AuthorizationRequest struct {
	MTI                   *field.String `index:"0"`
	PrimaryAccountNumber  *field.String `index:"2"`
	Amount                *field.String `index:"3"`
	TransmissionDateTime  *field.String `index:"4"`
	Currency              *field.String `index:"7"`
	CardVerificationValue *field.String `index:"8"`
	ExpirationDate        *field.String `index:"9"`
	STAN                  *field.String `index:"11"`
}

type AuthorizationResponse struct {
	MTI               *field.String `index:"0"`
	ApprovalCode      *field.String `index:"5"`
	AuthorizationCode *field.String `index:"6"`
	STAN              *field.String `index:"11"`
}
