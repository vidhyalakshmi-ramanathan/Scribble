package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/dto"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

var (
	ErrUserAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "User already exists",
	}
	ErrInvalidRequest = &AppError{
		Code:    http.StatusConflict,
		Message: "Invalid request",
	}
	ErrInvalidPassword = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid password",
	}
	ErrUnauthorizedAccess = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized access",
	}
	ErrInvalidToken = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Missing or Invalid token",
	}
	ErrDocumentNotExists = &AppError{
		Code:    http.StatusConflict,
		Message: "Document doesn't exists",
	}
	ErrOwnerNotExists = &AppError{
		Code:    http.StatusConflict,
		Message: "Owner doesn't exists",
	}
	ErrForbidden = &AppError{
		Code:    http.StatusForbidden,
		Message: "You do not have the access",
	}
	ErrUserAlreadyMember = &AppError{
		Code:    http.StatusForbidden,
		Message: "User is already a member",
	}
	ErrUserNotMember = &AppError{
		Code:    http.StatusForbidden,
		Message: "User is not a member",
	}
	ErrInvalidContent = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid edit content",
	}
	ErrVersionConflict = &AppError{
		Code:    http.StatusConflict,
		Message: "Document has changed. Please refresh or try again",
	}
	ErrInviteNotExists = &AppError{
		Code:    http.StatusConflict,
		Message: "Invite doesn't exists",
	}
	ErrInviteNotPendingOrExpired = &AppError{
		Code:    http.StatusConflict,
		Message: "Invite status either not pending or expired",
	}
	ErrCannotInviteSelf = &AppError{
		Code:    http.StatusBadRequest,
		Message: "You cannot invite yourself",
	}
	ErrInviteAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "Invite already exists",
	}
)

func (e *AppError) Error() string {
	return e.Message
}

var ErrNotFound = errors.New("resource not found")

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func Error(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "Something went wrong"
	var appErr *AppError
	var validateError *ValidationError

	log.Println("error: ", err)

	switch {
	case errors.Is(err, ErrNotFound):
		status = http.StatusNotFound
		message = "Resource not found"
	case errors.As(err, &appErr):
		status = appErr.Code
		message = appErr.Message
	case errors.As(err, &validateError):
		status = http.StatusBadRequest
		message = validateError.Message
	}
	JSON(w, status, dto.APIResponse{
		Message: message,
	})
}

type ValidationError struct {
	Message string
	Code    int
}

func (e *ValidationError) Error() string {
	return e.Message
}
