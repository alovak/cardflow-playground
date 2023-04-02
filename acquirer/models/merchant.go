package models

type CreateMerchant struct {
	Name       string
	MCC        string // Merchant Category Code
	PostalCode string
	WebSite    string
}

type Merchant struct {
	ID         string
	Name       string
	MCC        string // Merchant Category Code
	PostalCode string
	WebSite    string
}
