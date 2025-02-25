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

// CreateUser godoc
// @Summary Create a new user account
// @Description Creates a new user account and sends the request to Kafka for processing
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body models.Account true "Account details"
// @Success 200 {object} map[string]interface{} "success: true, msg: account creation request placed successfully"
// @Failure 400 {object} map[string]string "error: Invalid request body"
// @Failure 500 {object} map[string]string "error: Internal server error or Kafka failure"
// @Router /accounts [post]
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

// CreditAmount godoc
// @Summary Credit an amount to an account
// @Description Records a credit transaction and sends it to Kafka
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction details"
// @Success 200 {object} map[string]interface{} "success: true, msg: Credit Transaction Successfully Recorded"
// @Failure 400 {object} map[string]string "error: Invalid request body"
// @Failure 500 {object} map[string]string "error: Internal server error or Kafka failure"
// @Router /credit [post]
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

// WithdrawAmount godoc
// @Summary Withdraw an amount from an account
// @Description Records a withdrawal transaction and sends it to Kafka
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction details"
// @Success 200 {object} map[string]interface{} "success: true, msg: Withdraw Transaction Successfully Recorded"
// @Failure 400 {object} map[string]string "error: Invalid request body"
// @Failure 500 {object} map[string]string "error: Internal server error or Kafka failure"
// @Router /debit [post]
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

// TransferAmount godoc
// @Summary Transfer an amount between accounts
// @Description Records a transfer transaction and sends it to Kafka
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction details"
// @Success 200 {object} map[string]interface{} "success: true, msg: Transfer Transaction Successfully Recorded"
// @Failure 400 {object} map[string]string "error: Invalid request body"
// @Failure 500 {object} map[string]string "error: Internal server error or Kafka failure"
// @Router /transfer [post]
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

// FindTransactionHistory godoc
// @Summary Retrieve transaction history for an account
// @Description Fetches a list of transactions associated with the specified account number from the ledger.
// @Tags transactions
// @Accept json
// @Produce json
// @Param accountNumber path string true "Account Number" example:"Acc1234" description:"The unique identifier of the account (e.g., 'Acc1234')"
// @Success 200 {array} models.TransactionLedger "Successful response with a list of transactions"
// @Success 200 {object} []models.TransactionLedger "Empty list if no transactions are found"
// @Failure 400 {object} map[string]string "error: Account number is required" example:{"error":"Account number is required"}
// @Failure 500 {object} map[string]string "error: Failed to get transactions" example:{"error":"Failed to get transactions: database connection error"}
// @Router /transactions/{accountNumber} [get]
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
