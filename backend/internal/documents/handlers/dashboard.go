package handlers

import (
	"net/http"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/services"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/dto"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/middleware"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
)

type DashboardHandler struct {
	service services.DashboardService
}

func NewDashboardHandler(service services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dashboard, err := h.service.GetDashboard(claims.UserID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Dashboard Data Retrieved!",
		Data:    dashboard,
	})
}
