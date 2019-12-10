package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

// Response struct which contains an API Response
type Response struct {
	Message     string                 `json:"message,omitempty"`
	Validations []string               `json:"validations,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Successful  bool                   `json:"successful"`
}

// Transaction struct is a representation of Transaction fields
type Transaction struct {
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	StatusCode     int    `json:"statusCode"`
	Status         string `json:"status"`
	OrderReference string `json:"orderReference"`
	TransactionID  string `json:"transactionID"`
}

// FlypayA struct which contains an array of FlypayA Transactions
type FlypayA struct {
	Transactions []FlypayATransaction `json:"transactions"`
}

// FlypayATransaction struct is a representation of flypayA Transaction
type FlypayATransaction struct {
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	StatusCode     int    `json:"statusCode"`
	Status         string
	OrderReference string `json:"orderReference"`
	TransactionID  string `json:"transactionId"`
}

func (ft *FlypayATransaction) setStatus() {

	switch ft.StatusCode {
	case 1:
		ft.Status = "authorised"
	case 2:
		ft.Status = "decline"
	case 3:
		ft.Status = "refunded"
	}
}

// FlypayB struct which contains an array of FlypayB Transactions
type FlypayB struct {
	Transactions []FlypayBTransaction `json:"transactions"`
}

// FlypayBTransaction struct is a representation of flypayB Transaction
type FlypayBTransaction struct {
	Amount         int    `json:"value"`
	Currency       string `json:"transactionCurrency"`
	StatusCode     int    `json:"statusCode"`
	Status         string
	OrderReference string `json:"orderInfo"`
	TransactionID  string `json:"paymentId"`
}

func (ft *FlypayBTransaction) setStatus() {

	switch ft.StatusCode {
	case 100:
		ft.Status = "authorised"
	case 200:
		ft.Status = "decline"
	case 300:
		ft.Status = "refunded"
	}
}

func getFlypayA(ft TransactionRequestQuery) []Transaction {
	//TODO: handle large json files concurrently
	byteValue, _ := ioutil.ReadFile("flypayA.json")
	var flypayA FlypayA

	json.Unmarshal(byteValue, &flypayA)

	// Apply Query and Transform
	var data []Transaction

	for i := range flypayA.Transactions {
		selected := true
		flypayA.Transactions[i].setStatus()

		if ft.Status != "" && ft.Status != flypayA.Transactions[i].Status {
			selected = false
		}

		if ft.Currency != "" && ft.Currency != flypayA.Transactions[i].Currency {
			selected = false
		}

		fmt.Println(ft.AmountMin > 0 && flypayA.Transactions[i].Amount < ft.AmountMin, ft.AmountMin, flypayA.Transactions[i].Amount)
		if ft.AmountMin > 0 && flypayA.Transactions[i].Amount < ft.AmountMin {
			selected = false
		}

		if ft.AmountMax > 0 && flypayA.Transactions[i].Amount > ft.AmountMax {
			selected = false
		}

		if selected == true {
			data = append(data, Transaction(flypayA.Transactions[i]))
		}

	}
	return data
}

func getFlypayB(ft TransactionRequestQuery) []Transaction {
	//TODO: handle large json files concurrently
	byteValue, _ := ioutil.ReadFile("flypayB.json")
	var flypayB FlypayB

	json.Unmarshal(byteValue, &flypayB)

	// Apply Query and Transform
	var data []Transaction

	for i := range flypayB.Transactions {
		selected := true
		flypayB.Transactions[i].setStatus()

		if ft.Status != "" && ft.Status != flypayB.Transactions[i].Status {
			selected = false
		}

		if ft.Currency != "" && ft.Currency != flypayB.Transactions[i].Currency {
			selected = false
		}

		if ft.AmountMin > 0 && flypayB.Transactions[i].Amount < ft.AmountMin {
			selected = false
		}

		if ft.AmountMax > 0 && flypayB.Transactions[i].Amount > ft.AmountMax {
			selected = false
		}

		if selected == true {
			data = append(data, Transaction(flypayB.Transactions[i]))
		}

	}

	return data
}

var wg = sync.WaitGroup{}

func findTransactions(ft TransactionRequestQuery) ([]Transaction, error) {
	data := make([]Transaction, 0)

	if ft.Provider == "" || (ft.Provider != "" && ft.Provider == "flypayA") {
		wg.Add(1)
		go func(ft TransactionRequestQuery) {
			transactions := getFlypayA(ft)
			fmt.Println(transactions)
			data = append(data, transactions...)
			wg.Done()
		}(ft)
	}

	if ft.Provider == "" || (ft.Provider != "" && ft.Provider == "flypayB") {
		wg.Add(1)
		go func(ft TransactionRequestQuery) {
			transactions := getFlypayB(ft)
			fmt.Println(transactions)
			data = append(data, transactions...)
			wg.Done()
		}(ft)
	}

	wg.Wait()
	return data, nil
}

// TransactionRequestQuery struct is a representation of Transaction Request Query
type TransactionRequestQuery struct {
	Provider  string `json:"provider,omitempty"`
	Status    string `json:"status,omitempty"`
	Currency  string `json:"currency,omitempty"`
	AmountMin int    `json:"amountMin,omitempty"`
	AmountMax int    `json:"amountMax,omitempty"`
}

func (a *TransactionRequestQuery) validate() []string {
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

func findTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	var transactionRequestQuery TransactionRequestQuery
	decoder := schema.NewDecoder()

	err := decoder.Decode(&transactionRequestQuery, r.URL.Query())
	if err != nil {
		//TODO: Send Not Allowed Fields response
		fmt.Println(err)
	}

	validErrs := transactionRequestQuery.validate()
	if len(validErrs) > 0 {
		response := Response{Message: "Validations Errors", Validations: validErrs, Successful: false}

		payload, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}

	data, err := findTransactions(transactionRequestQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := Response{Message: "Data retrieved successfully", Data: map[string]interface{}{"transactions": data}, Successful: true}

	payload, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Not found", Successful: false}

	payload, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func main() {

	// Init Router
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	// Route Handlers - Endpoints
	r.HandleFunc("/api/payment/transaction", findTransactionsHandler).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(notFound)

	log.Fatal(http.ListenAndServe(":3000", r))
}
