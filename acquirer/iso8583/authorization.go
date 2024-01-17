package iso8583

type AuthorizationRequest struct {
	MTI                   string               `index:"0"`
	PrimaryAccountNumber  string               `index:"2"`
	Amount                int64                `index:"3"`
	TransmissionDateTime  string               `index:"4"`
	Currency              string               `index:"7"`
	CardVerificationValue string               `index:"8"`
	ExpirationDate        string               `index:"9"`
	AcceptorInformation   *AcceptorInformation `index:"10"`
	STAN                  string               `index:"11"`
}

type AuthorizationResponse struct {
	MTI               string `index:"0"`
	ApprovalCode      string `index:"5"`
	AuthorizationCode string `index:"6"`
	STAN              string `index:"11"`
}

type AcceptorInformation struct {
	Name       string `index:"01"`
	MCC        string `index:"02"`
	PostalCode string `index:"03"`
	WebSite    string `index:"04"`
}
