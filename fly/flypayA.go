package fly

import (
	"log"
	"os"
)

// FlypayATransaction struct is a representation of flypayA Transaction
type flypayATransaction struct {
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	StatusCode     int    `json:"statusCode"`
	Status         string
	OrderReference string `json:"orderReference"`
	TransactionID  string `json:"transactionId"`
}

func (ft *flypayATransaction) setStatus() {

	switch ft.StatusCode {
	case 1:
		ft.Status = "authorised"
	case 2:
		ft.Status = "decline"
	case 3:
		ft.Status = "refunded"
	}
}

// GetFlypayA return FlypayA Transactions
func GetFlypayA(ft TransactionRequestQuery) []Transaction {
	// Open the file
	input, err := os.Open("flypayA.json")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan flypayATransaction)
	go streamAParsing(input, c)

	// Apply Query and Transform
	var data []Transaction

	for {
		tr, open := <-c

		if !open {
			break
		}

		selected := true
		tr.setStatus()

		if ft.Status != "" && ft.Status != tr.Status {
			selected = false
		}

		if ft.Currency != "" && ft.Currency != tr.Currency {
			selected = false
		}

		if ft.AmountMin > 0 && tr.Amount < ft.AmountMin {
			selected = false
		}

		if ft.AmountMax > 0 && tr.Amount > ft.AmountMax {
			selected = false
		}

		if selected == true {
			data = append(data, Transaction(tr))
		}
	}

	return data
}
