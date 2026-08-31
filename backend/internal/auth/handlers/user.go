package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/entities/packets"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/services"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/dto"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) SignUpUser(w http.ResponseWriter, r *http.Request) {
	var userData packets.UserDataReq

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "The JSON in the request is invalid", http.StatusBadRequest)
		return
	}

	user := entities.User{
		UserID:   uuid.New().String(),
		UserName: userData.Username,
		Email:    userData.Email,
		Password: userData.Password,
	}
	if err := h.service.SignUpUser(&user); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, dto.APIResponse{
		Message: "User SignUp Successfully!",
	})
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginData packets.LoginReq

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "The JSON in the request is invalid", http.StatusBadRequest)
		return
	}

	user := packets.LoginReq{
		Email:    loginData.Email,
		Password: loginData.Password,
	}
	loginResponse, err := h.service.Login(&user)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "User Login Successfully!",
		Data:    loginResponse,
	})
}

func (h *UserHandler) Test(w http.ResponseWriter, r *http.Request) {

	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Tested Successfully!",
	})
}
