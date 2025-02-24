package handlers

import (
	"accountProducer/kafka"
	"accountProducer/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	kafkactl *kafka.KafkaController
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler() *AccountHandler {
	return &AccountHandler{
		kafkactl: &kafka.KafkaController{},
	}
}

func (h *AccountHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// decoding
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

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

func (h *AccountHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/accounts", h.CreateUser).Methods("POST")
}
