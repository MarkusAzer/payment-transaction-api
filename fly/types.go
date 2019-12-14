package fly

// Transaction struct is a representation of Transaction fields
type Transaction struct {
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	StatusCode     int    `json:"statusCode"`
	Status         string `json:"status"`
	OrderReference string `json:"orderReference"`
	TransactionID  string `json:"transactionID"`
}

// TransactionRequestQuery struct is a representation of Transaction Request Query
type TransactionRequestQuery struct {
	Provider  string `json:"provider,omitempty"`
	Status    string `json:"status,omitempty"`
	Currency  string `json:"currency,omitempty"`
	AmountMin int    `json:"amountMin,omitempty"`
	AmountMax int    `json:"amountMax,omitempty"`
}

// Validate Method is to validate Transaction Request Query
func (a *TransactionRequestQuery) Validate() []string {
	var errs []string

	if a.Provider != "" && a.Provider != "flypayA" && a.Provider != "flypayB" {
		errs = append(errs, "Provider :- Valid Providers flypayA, flypayB")
	}

	if a.Status != "" && a.Status != "authorised" && a.Status != "decline" && a.Status != "refunded" {
		errs = append(errs, "Status :- Valid Status authorised, decline, refunded")
	}

	if (a.AmountMin != 0) && (a.AmountMin < 1) {
		errs = append(errs, "AmountMin :- Must be greater than or equal 1")
	}

	if (a.AmountMax != 0) && a.AmountMax < 1 {
		errs = append(errs, "AmountMax :- Must be greater than or equal 1")
	}

	return errs
}
