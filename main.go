package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/markusazer/payment-transaction-api/fly"
)

// Response struct which contains an API Response
type Response struct {
	Message     string                 `json:"message,omitempty"`
	Validations []string               `json:"validations,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Successful  bool                   `json:"successful"`
}

func findTransactions(ft fly.TransactionRequestQuery) ([]fly.Transaction, error) {
	var wg = sync.WaitGroup{}
	data := make([]fly.Transaction, 0)

	if ft.Provider == "" || (ft.Provider != "" && ft.Provider == "flypayA") {
		wg.Add(1)
		go func(ft fly.TransactionRequestQuery) {
			transactions := fly.GetFlypayA(ft)
			data = append(data, transactions...)
			wg.Done()
		}(ft)
	}

	if ft.Provider == "" || (ft.Provider != "" && ft.Provider == "flypayB") {
		wg.Add(1)
		go func(ft fly.TransactionRequestQuery) {
			transactions := fly.GetFlypayB(ft)
			data = append(data, transactions...)
			wg.Done()
		}(ft)
	}

	wg.Wait()
	return data, nil
}

func findTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	// Get Query Params
	var transactionRequestQuery fly.TransactionRequestQuery
	decoder := schema.NewDecoder()

	err := decoder.Decode(&transactionRequestQuery, r.URL.Query())
	if err != nil {
		//TODO: Send Not Allowed Fields response
		fmt.Println(err)
	}

	// Validate Query Params
	validErrs := transactionRequestQuery.Validate()
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
		return
	}

	// Get transactions
	data, err := findTransactions(transactionRequestQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create and Send Response
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
