package handlers

import (
	"accountProducer/database"
	"accountProducer/kafka"
	"accountProducer/models"
	"accountProducer/repositories"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type AccountHandler struct {
	kafkactl *kafka.KafkaController
	trrepo   repositories.Repository
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(db database.Database, lobbs *hclog.Logger) *AccountHandler {
	return &AccountHandler{
		kafkactl: &kafka.KafkaController{},
		trrepo:   repositories.NewTransactionRepository(db, lobbs),
	}
}

func (h *AccountHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// decoding
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println(account)

	accountInBytes, err := json.Marshal(account)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
		return
	}

	// send the bytes to kafka
	err = h.kafkactl.PushOrderToQueue("account-creation", accountInBytes)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Not able to send message to kafka", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"msg":     "account creation request placed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println(err)
		http.Error(w, "Kya hua pta nhi", http.StatusInternalServerError)
		return
	}
}

func (h *AccountHandler) CreditAmount(w http.ResponseWriter, r *http.Request) {
	// decoding
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transactionInBytes, err := json.Marshal(transaction)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
		return
	}

	// send the bytes to kafka
	err = h.kafkactl.PushOrderToQueue("transaction", transactionInBytes)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Not able to send message to kafka", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"msg":     "Credit Transaction Successfully Recorded",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println(err)
		http.Error(w, "Encode Not Happen Properly", http.StatusInternalServerError)
		return
	}
}

func (h *AccountHandler) WithdrawAmount(w http.ResponseWriter, r *http.Request) {
	// decoding
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transactionInBytes, err := json.Marshal(transaction)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
		return
	}

	// send the bytes to kafka
	err = h.kafkactl.PushOrderToQueue("transaction", transactionInBytes)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Not able to send message to kafka", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"msg":     "Withdraw Transaction Successfully Recorded",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println(err)
		http.Error(w, "Encode Not Happen Properly", http.StatusInternalServerError)
		return
	}
}

func (h *AccountHandler) TransferAmount(w http.ResponseWriter, r *http.Request) {
	// decoding
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transactionInBytes, err := json.Marshal(transaction)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
		return
	}

	// send the bytes to kafka
	err = h.kafkactl.PushOrderToQueue("transaction", transactionInBytes)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Not able to send message to kafka", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"msg":     "Transfer Transaction Successfully Recorded",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println(err)
		http.Error(w, "Encode Not Happen Properly", http.StatusInternalServerError)
		return
	}
}

// Handler function
func (h *AccountHandler) FindTransactionHistory(w http.ResponseWriter, r *http.Request) {
	// Extract path variables
	vars := mux.Vars(r)
	accountNumber := vars["accountNumber"] // Get the value of {accountNumber}

	if accountNumber == "" {
		http.Error(w, "Account number is required", http.StatusBadRequest)
		return
	}

	// Query transactions
	transactions, err := h.trrepo.FindTransactionByAccountNumber(r.Context(), accountNumber)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *AccountHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/accounts", h.CreateUser).Methods("POST")
	router.HandleFunc("/debit", h.WithdrawAmount).Methods("POST")
	router.HandleFunc("/credit", h.CreditAmount).Methods("POST")
	router.HandleFunc("/transfer", h.TransferAmount).Methods("POST")
	router.HandleFunc("/transactions/{accountNumber}", h.FindTransactionHistory).Methods("GET")
}
