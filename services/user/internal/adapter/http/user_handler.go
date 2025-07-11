package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robrt95x/godops/pkg/errors"
	"github.com/robrt95x/godops/services/user/internal/application/usecase"
)

type UserHandler struct {
	createUserUseCase *usecase.CreateUserUseCase
	getUserUseCase    *usecase.GetUserUseCase
}

func NewUserHandler(createUserUseCase *usecase.CreateUserUseCase, getUserUseCase *usecase.GetUserUseCase) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUseCase,
		getUserUseCase:    getUserUseCase,
	}
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.BadRequest(w, err)
		return
	}
	
	user, err := h.createUserUseCase.Execute(req.Name, req.Email)
	if err != nil {
		errors.BadRequest(w, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	user, err := h.getUserUseCase.Execute(id)
	if err != nil {
		errors.NotFound(w, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.getUserUseCase.ExecuteGetAll()
	if err != nil {
		errors.InternalServerError(w, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
