package fly

import (
	"log"
	"os"
)

// FlypayBTransaction struct is a representation of flypayB Transaction
type flypayBTransaction struct {
	Amount         int    `json:"value"`
	Currency       string `json:"transactionCurrency"`
	StatusCode     int    `json:"statusCode"`
	Status         string
	OrderReference string `json:"orderInfo"`
	TransactionID  string `json:"paymentId"`
}

func (ft *flypayBTransaction) setStatus() {

	switch ft.StatusCode {
	case 100:
		ft.Status = "authorised"
	case 200:
		ft.Status = "decline"
	case 300:
		ft.Status = "refunded"
	}
}

// GetFlypayB return FlypayB Transactions
func GetFlypayB(ft TransactionRequestQuery) []Transaction {
	// Open the file
	input, err := os.Open("flypayB.json")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan flypayBTransaction)
	go streamBParsing(input, c)

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
