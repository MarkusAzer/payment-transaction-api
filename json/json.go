package json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// FlypayATransaction struct is a representation of flypayA Transaction
type FlypayATransaction struct {
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	StatusCode     int    `json:"statusCode"`
	Status         string
	OrderReference string `json:"orderReference"`
	TransactionID  string `json:"transactionId"`
}

func StreamParsing() {

	// Error Handling
	he := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	// Open the file
	input, err := os.Open("flypayA.json")
	he(err)

	// Start reading the file
	dec := json.NewDecoder(bufio.NewReader(input))

	// Start of the expect object
	t, err := dec.Token()
	he(err)

	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		log.Fatal("Expected object")
	}

	// Read props
	for dec.More() {
		t, err = dec.Token()
		he(err)

		// transactions array
		if t == "transactions" {
			t, err := dec.Token()
			he(err)
			if delim, ok := t.(json.Delim); !ok || delim != '[' {
				log.Fatal("Expected array")
			}

			// Read transactions
			for dec.More() {
				// Read next transaction
				ft := FlypayATransaction{}
				he(dec.Decode(&ft))
				fmt.Printf("transaction: %+v\n", ft)
			}

			// Array closing delim
			t, err = dec.Token()
			he(err)
			if delim, ok := t.(json.Delim); !ok || delim != ']' {
				log.Fatal("Expected array closing")
			}
		}
	}

	// Object closing delim
	t, err = dec.Token()
	he(err)
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		log.Fatal("Expected object closing")
	}
}
