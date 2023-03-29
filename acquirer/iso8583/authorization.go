package iso8583

import "github.com/moov-io/iso8583/field"

type AuthorizationRequest struct {
	MTI                  *field.String `iso8583:"0"`
	PrimaryAccountNumber *field.String `iso8583:"2"`
	Amount               *field.String `iso8583:"3"`
	TransmissionDateTime *field.String `iso8583:"4"`
	STAN                 *field.String `iso8583:"11"`
}

type AuthorizationResponse struct {
	MTI               *field.String `iso8583:"0"`
	ApprovalCode      *field.String `iso8583:"5"`
	AuthorizationCode *field.String `iso8583:"6"`
	STAN              *field.String `iso8583:"11"`
}
