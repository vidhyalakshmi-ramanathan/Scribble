package collaboration

import (
	"encoding/json"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/services"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"gorm.io/datatypes"
)

type CollaborationService struct {
	documentService services.DocumentService
}

func NewCollaborationService(documentService services.DocumentService) *CollaborationService {
	return &CollaborationService{documentService: documentService}
}
func (s *CollaborationService) ApplyEdit(userID, documentID string, editData entities.EditData) (*entities.Documents, error) {
	// Validate the submitted TipTap JSON.
	if len(editData.Content) == 0 || !json.Valid(editData.Content) {
		return nil, response.ErrInvalidContent
	}

	var content struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(editData.Content, &content); err != nil || content.Type != "doc" {
		return nil, response.ErrInvalidContent
	}

	document, err := s.documentService.CheckEditPermission(userID, documentID)
	if err != nil {
		return nil, err
	}

	if editData.BaseCollabVersion != document.CollabVersion {
		return nil, response.ErrVersionConflict
	}

	// Content is a complete TipTap JSON document
	document.Content = datatypes.JSON(editData.Content)
	document.UpdatedBy = userID

	// Atomic content update and concurrency-version increment
	if err := s.documentService.ApplyCollaborativeEdit(document, editData.BaseCollabVersion); err != nil {
		return nil, err
	}

	document.CollabVersion = editData.BaseCollabVersion + 1
	return document, nil
}
