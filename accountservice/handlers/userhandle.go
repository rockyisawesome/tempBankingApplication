package handlers

import (
	"accountservice/database"
	"accountservice/models"
	"accountservice/repositories"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	repo repositories.Repository
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(db *database.PostgresPoolDB) *UserHandler {
	return &UserHandler{
		repo: repositories.NewUserRepository(db),
	}
}

// CreateUser handles POST requests to create a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// decoding
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now
	account.IsActive = true

	ctx := r.Context()
	if err := h.repo.Create(ctx, &account); err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// GetUser handles GET requests to retrieve a user by ID
// func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	idStr := vars["id"]
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	ctx := r.Context()
// 	account, err := h.repo.GetByID(ctx, id)
// 	if err != nil {
// 		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if account == nil {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(account)
// }

// registering the routes

func (h *UserHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/accounts", h.CreateUser).Methods("POST")
}
