package models

type Merchant struct {
	Name       string
	MCC        string // Merchant Category Code
	PostalCode string
	WebSite    string
}
