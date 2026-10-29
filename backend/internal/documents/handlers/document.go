package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/services"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/dto"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/middleware"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"gorm.io/datatypes"
)

type DocumentHandler struct {
	service services.DocumentService
}

func NewDocumentHandler(service services.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

func (h *DocumentHandler) CreateDoc(w http.ResponseWriter, r *http.Request) {
	var docData entities.DocumentReq
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&docData); err != nil {
		http.Error(w, "The JSON in the request is invalid", http.StatusBadRequest)
		return
	}

	document := entities.Documents{
		DocumentID:    uuid.New().String(),
		Title:         docData.Title,
		Content:       datatypes.JSON(docData.Content),
		Category:      docData.Category,
		OwnerID:       claims.UserID,
		UpdatedBy:     claims.UserID,
		SchemaVer:     1,
		Version:       0,
		CollabVersion: 0,
	}
	docID, err := h.service.CreateDoc(claims.UserID, &document)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, dto.APIResponse{
		Message: "Document Created Successfully!",
		Data:    map[string]string{"docID": docID},
	})
}

func (h *DocumentHandler) GetDocByID(w http.ResponseWriter, r *http.Request) {
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
	docDetails, err := h.service.GetDocByID(claims.UserID, docID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Document Retrieved Successfully!",
		Data:    docDetails,
	})
}

func (h *DocumentHandler) GetMyAllDocs(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	docDetails, err := h.service.GetMyAllDocs(claims.UserID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Documents list retrieved for a user!",
		Data:    docDetails,
	})
}

func (h *DocumentHandler) UpdateDoc(w http.ResponseWriter, r *http.Request) {
	var docData entities.DocumentReq
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&docData); err != nil {
		http.Error(w, "The JSON in the request is invalid", http.StatusBadRequest)
		return
	}
	docID := chi.URLParam(r, "docID")
	if docID == "" {
		http.Error(w, "Request parameter docID is required", http.StatusBadRequest)
		return
	}
	if docData.CollabVersion == nil || *docData.CollabVersion < 0 {
		http.Error(w, "A valid collaboration version is required", http.StatusBadRequest)
		return
	}
	document := entities.Documents{
		DocumentID:    docID,
		Title:         docData.Title,
		Content:       datatypes.JSON(docData.Content),
		Category:      docData.Category,
		UpdatedBy:     claims.UserID,
		CollabVersion: *docData.CollabVersion,
	}
	if err := h.service.UpdateDoc(claims.UserID, &document); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Document Updated Successfully!",
	})
}

func (h *DocumentHandler) DeleteDoc(w http.ResponseWriter, r *http.Request) {
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
	if err := h.service.DeleteDoc(claims.UserID, docID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Document Deleted Successfully!",
	})
}

func (h *DocumentHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
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
	versionParam := chi.URLParam(r, "version")
	if versionParam == "" {
		http.Error(w, "Request parameter version is required", http.StatusBadRequest)
		return
	}
	version, err := strconv.Atoi(versionParam)
	if err != nil || version < 1 {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}
	versionDetail, err := h.service.GetVersion(claims.UserID, docID, version)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Retrieved document version!",
		Data:    versionDetail,
	})
}

func (h *DocumentHandler) GetVersions(w http.ResponseWriter, r *http.Request) {
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
	versions, err := h.service.GetVersions(claims.UserID, docID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Retrieved document history!",
		Data:    versions,
	})
}

func (h *DocumentHandler) RestoreVersion(w http.ResponseWriter, r *http.Request) {
	var req entities.RestoreVersionReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	if req.ExpectedCollabVersion == nil || *req.ExpectedCollabVersion < 0 {
		http.Error(w, "A valid expectedCollabVersion is required", http.StatusBadRequest)
		return
	}

	docID := chi.URLParam(r, "docID")
	if docID == "" {
		http.Error(w, "Request parameter docID is required", http.StatusBadRequest)
		return
	}
	versionParam := chi.URLParam(r, "version")
	if versionParam == "" {
		http.Error(w, "Request parameter version is required", http.StatusBadRequest)
		return
	}
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	version, err := strconv.Atoi(versionParam)
	if err != nil || version < 1 {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}
	err = h.service.RestoreVersion(claims.UserID, docID, version, *req.ExpectedCollabVersion)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Document restored successfully!",
	})
}

func (h *DocumentHandler) SaveVersion(w http.ResponseWriter, r *http.Request) {
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

	err := h.service.SaveVersion(claims.UserID, docID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.APIResponse{
		Message: "Document Version successfully!",
	})
}
