package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/services"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/dto"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/middleware"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
)

type RoleHandler struct {
	service services.RoleService
}

func NewRoleHandler(service services.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) GetMemberDetails(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "docID")
	if docID == "" {
		http.Error(w, "Request parameter docID is required", http.StatusBadRequest)
		return
	}
	memberID := chi.URLParam(r, "memberID")
	if memberID == "" {
		http.Error(w, "Request parameter userID is required", http.StatusBadRequest)
		return
	}
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	member, err := h.service.GetMemberDetails(claims.UserID, memberID, docID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Member Detail Retrieved!",
		Data:    member,
	})
}

func (h *RoleHandler) GetMembersList(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "docID")
	if docID == "" {
		http.Error(w, "Request parameter docID is required", http.StatusBadRequest)
		return
	}
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	membersList, err := h.service.GetMembersList(claims.UserID, docID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Members list retrieved!",
		Data:    membersList,
	})
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "docID")
	if docID == "" {
		http.Error(w, "Request parameter docID is required", http.StatusBadRequest)
		return
	}
	memberID := chi.URLParam(r, "memberID")
	if memberID == "" {
		http.Error(w, "Request parameter userID is required", http.StatusBadRequest)
		return
	}
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var role entities.UpdateRoleReq
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, "The JSON in the request is invalid", http.StatusBadRequest)
		return
	}
	err := h.service.UpdateRole(claims.UserID, memberID, docID, role.Role)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, dto.APIResponse{
		Message: "Member role updated successfully!",
	})
}

func (h *RoleHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "docID")
	if docID == "" {
		http.Error(w, "Request parameter docID is required", http.StatusBadRequest)
		return
	}
	memberID := chi.URLParam(r, "memberID")
	if memberID == "" {
		http.Error(w, "Request parameter userID is required", http.StatusBadRequest)
		return
	}
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err := h.service.DeleteMember(claims.UserID, memberID, docID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Member Deleted Successfully!",
	})
}

func (h *RoleHandler) CreateInvite(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "docID")
	if documentID == "" {
		http.Error(w, "Request parameter docID is required", http.StatusBadRequest)
		return
	}
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req *entities.AddUserRoleReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "The JSON in the request is invalid", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateInvite(claims.UserID, documentID, req.Email, req.Role); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, dto.APIResponse{
		Message: "Invited Member Successfully!",
	})
}

func (h *RoleHandler) GetSentInvites(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "docID")
	if documentID == "" {
		http.Error(w, "Request parameter docID is required", http.StatusBadRequest)
		return
	}
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	invites, err := h.service.GetSentInvites(claims.UserID, documentID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Sent Invites Retrieved!",
		Data:    invites,
	})
}

func (h *RoleHandler) GetInvites(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	invites, err := h.service.GetMyInvites(claims.UserID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Invites Retrieved!",
		Data:    invites,
	})
}

func (h *RoleHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	inviteID := chi.URLParam(r, "inviteID")
	if inviteID == "" {
		http.Error(w, "Request parameter inviteID is required", http.StatusBadRequest)
		return
	}
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := h.service.AcceptInvite(claims.UserID, inviteID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Member Added Successfully!",
	})
}
