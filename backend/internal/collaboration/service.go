package collaboration

import (
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/services"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
)

type CollaborationService struct {
	documentService services.DocumentService
}

func NewCollaborationService(documentService services.DocumentService) *CollaborationService {
	return &CollaborationService{documentService: documentService}
}
func (s *CollaborationService) ApplyEdit(userID, documentID string, editData entities.EditData) (*entities.Documents, error) {
	document, err := s.documentService.CheckEditPermission(userID, documentID)
	if err != nil {
		return nil, err
	}

	if editData.BaseVersion != document.Version {
		return nil, response.ErrVersionConflict
	}

	switch editData.Operation {
	case "insert":
		if editData.Position < 0 || editData.Position > len(document.Content) {
			return nil, response.ErrInvalidPosition
		}
		document.Content = document.Content[:editData.Position] + editData.Text + document.Content[editData.Position:]
	case "delete":
		if editData.Position < 0 || editData.Position > len(document.Content) || editData.Length < 0 || editData.Length > len(document.Content)-editData.Position {
			return nil, response.ErrInvalidPosition
		}
		document.Content = document.Content[:editData.Position] + document.Content[editData.Position+editData.Length:]
	default:
		return nil, response.ErrInvalidOperation
	}

	// Atomic content update and concurrency-version increment
	document.UpdatedBy = userID

	if err := s.documentService.UpdateContentByVersion(document, editData.BaseVersion); err != nil {
		return nil, err
	}

	document.Version = editData.BaseVersion + 1
	return document, nil
}
