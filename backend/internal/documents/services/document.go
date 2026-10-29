package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	auth "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/models"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/models"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/enums"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DocumentService interface {
	CreateDoc(userID string, document *entities.Documents) (string, error)
	GetDocByID(userID, documentID string) (*entities.DocumentDetails, error)
	GetMyAllDocs(userID string) ([]entities.DocumentResponse, error)
	UpdateDoc(userID string, document *entities.Documents) error
	DeleteDoc(userID, documentID string) error
	CheckEditPermission(userID, docID string) (*entities.Documents, error)
	ApplyCollaborativeEdit(document *entities.Documents, baseVersion int) error
	GetVersion(userID, documentID string, version int) (versionDetail *entities.Version, err error)
	GetVersions(userID, documentID string) (versions []entities.VersionsList, err error)
	RestoreVersion(userID, documentID string, version, expectedCollabVersion int) error
	SaveVersion(userID, documentID string) error
}

type documentService struct {
	userModel    auth.UserModel
	model        models.DocumentModel
	rolemodel    models.RoleModel
	versionModel models.VersionModel
	db           *gorm.DB
}

func NewDocumentService(userModel auth.UserModel, model models.DocumentModel, rolemodel models.RoleModel, versionModel models.VersionModel, db *gorm.DB) DocumentService {
	return &documentService{
		userModel:    userModel,
		model:        model,
		rolemodel:    rolemodel,
		versionModel: versionModel,
		db:           db,
	}
}

func (s *documentService) CreateDoc(userID string, document *entities.Documents) (string, error) {
	if strings.TrimSpace(document.Title) == "" {
		document.Title = "Untitled Document"
	}
	if len(document.Content) == 0 {
		document.Content = datatypes.JSON(
			[]byte(`{"type":"doc","content":[{"type":"paragraph"}]}`),
		)
	}
	if !json.Valid(document.Content) {
		return "", &response.ValidationError{
			Message: "content must be valid JSON",
		}
	}
	var content struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(document.Content, &content); err != nil || content.Type != "doc" {
		return "", &response.ValidationError{
			Message: "content must be a valid ProseMirror document",
		}
	}
	if document.Category == "" {
		document.Category = "General"
	}

	document.OwnerID = userID
	document.UpdatedBy = userID
	document.SchemaVer = 1
	document.Version = 0
	document.CollabVersion = 0

	if err := s.model.CreateDoc(document); err != nil {
		return "", err
	}
	return document.DocumentID, nil
}

func (s *documentService) GetDocByID(userID, documentID string) (*entities.DocumentDetails, error) {
	doc, err := s.CheckAccess(userID, documentID)
	if err != nil {
		return nil, err
	}
	var deletedAt *time.Time

	if doc.DeletedAt.Valid {
		deletedAt = &doc.DeletedAt.Time
	}
	return &entities.DocumentDetails{
		DocumentID:    doc.DocumentID,
		Title:         doc.Title,
		Content:       doc.Content,
		CollabVersion: doc.CollabVersion,
		OwnerID:       doc.OwnerID,
		Version:       doc.Version,
		Category:      doc.Category,
		CreatedAt:     doc.CreatedAt,
		UpdatedAt:     doc.UpdatedAt,
		DeletedAt:     deletedAt,
		UpdatedBy:     doc.UpdatedBy,
	}, nil
}

func (s *documentService) GetMyAllDocs(userID string) ([]entities.DocumentResponse, error) {
	documents, err := s.model.GetMyAllDocs(userID)
	if err != nil {
		return nil, err
	}
	result := make([]entities.DocumentResponse, 0, len(documents))

	for _, doc := range documents {
		result = append(result, entities.DocumentResponse{
			DocumentID:    doc.DocumentID,
			Title:         doc.Title,
			CollabVersion: doc.CollabVersion,
			OwnerID:       doc.OwnerID,
			Version:       doc.Version,
			Category:      doc.Category,
			UpdatedAt:     doc.UpdatedAt,
			UpdatedBy:     doc.UpdatedBy,
		})
	}
	return result, nil
}

// UpdateDoc When AutoSave
func (s *documentService) UpdateDoc(userID string, document *entities.Documents) error {
	if strings.TrimSpace(document.Title) == "" {
		return &response.ValidationError{
			Message: "title cannot be empty",
		}
	}
	if len(document.Content) == 0 || !json.Valid(document.Content) {
		return &response.ValidationError{
			Message: "content must be valid JSON",
		}
	}
	var content struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(document.Content, &content); err != nil || content.Type != "doc" {
		return &response.ValidationError{
			Message: "content must be a valid ProseMirror document",
		}
	}
	expectedCollabVersion := document.CollabVersion

	oldDocument, err := s.CheckEditPermission(userID, document.DocumentID)
	if err != nil {
		return err
	}
	// Reject a request based on an outdated document revision
	if oldDocument.CollabVersion != expectedCollabVersion {
		return response.ErrVersionConflict
	}

	oldDocument.Title = document.Title
	oldDocument.Content = document.Content
	oldDocument.Category = document.Category
	oldDocument.UpdatedBy = userID

	updated, err := s.model.UpdateDocWithVersion(oldDocument, expectedCollabVersion)
	if err != nil {
		return err
	}
	if !updated {
		return response.ErrVersionConflict
	}
	return nil
}

func (s *documentService) DeleteDoc(userID, documentID string) error {
	document, err := s.getOwnedDocument(documentID)
	if err != nil {
		return err
	}
	if document.OwnerID != userID {
		return response.ErrForbidden
	}
	return s.model.DeleteDoc(documentID)
}

func (s *documentService) getOwnedDocument(docID string) (document *entities.Documents, err error) {
	document, err = s.model.GetDocByID(docID)
	if err != nil {
		if errors.Is(err, response.ErrNotFound) {
			return nil, response.ErrDocumentNotExists
		}
		return nil, err
	}
	return document, nil
}

func (s *documentService) CheckAccess(userID, docID string) (*entities.Documents, error) {
	document, err := s.model.GetDocByID(docID)
	if err != nil {
		if errors.Is(err, response.ErrNotFound) {
			return nil, response.ErrDocumentNotExists
		}
		return nil, err
	}
	if document.OwnerID == userID {
		return document, nil
	}

	_, err = s.rolemodel.GetMemberDetails(userID, docID)
	if err != nil {
		return nil, response.ErrUserNotMember
	}
	return document, nil
}

func (s *documentService) CheckEditPermission(userID, docID string) (*entities.Documents, error) {
	document, err := s.model.GetDocByID(docID)
	if err != nil {
		return nil, err
	}
	if err == response.ErrNotFound {
		return nil, response.ErrDocumentNotExists
	}
	if document.OwnerID == userID {
		return document, nil
	}

	member, err := s.rolemodel.GetMemberDetails(userID, docID)
	if err != nil {
		return nil, response.ErrUserNotMember
	}
	if member.Role == enums.Admin || member.Role == enums.Editor {
		return document, nil
	}
	return nil, response.ErrForbidden
}

func (s *documentService) ApplyCollaborativeEdit(document *entities.Documents, baseCollabVersion int) error {

	updated, err := s.versionModel.ApplyCollaborativeEdit(document.DocumentID, document, baseCollabVersion)
	if err != nil {
		return err
	}
	if !updated {
		return response.ErrVersionConflict
	}
	return nil
}

func (s *documentService) GetVersion(userID, documentID string, version int) (versionDetail *entities.Version, err error) {
	_, err = s.CheckAccess(userID, documentID)
	if err != nil {
		return nil, err
	}
	versionDetail, err = s.versionModel.GetVersion(documentID, version)
	if err != nil {
		return nil, err
	}
	return versionDetail, nil
}

func (s *documentService) GetVersions(userID, documentID string) (versionsList []entities.VersionsList, err error) {
	_, err = s.CheckAccess(userID, documentID)
	if err != nil {
		return nil, err
	}
	versionsList, err = s.versionModel.GetVersions(documentID)
	if err != nil {
		return nil, err
	}
	return versionsList, nil
}

func (s *documentService) RestoreVersion(userID, documentID string, version, expectedCollabVersion int) error {
	if _, err := s.CheckEditPermission(userID, documentID); err != nil {
		return err
	}
	// Retrieve the snapshot the user wants to restore
	versionToRestore, err := s.versionModel.GetVersion(documentID, version)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Lock the document row while restoring
		document, err := s.model.GetDocumentForUpdate(tx, documentID)
		if err != nil {
			return err
		}
		if document.CollabVersion != expectedCollabVersion {
			return response.ErrVersionConflict
		}
		// Find the latest saved history snapshot
		latestVersion, err := s.versionModel.GetLatestVersion(tx, documentID)
		if err != nil {
			return err
		}
		nextVersion := 1
		sameAsLatest := false

		if latestVersion != nil {
			nextVersion = latestVersion.Version + 1
			sameAsLatest = latestVersion.Title == versionToRestore.Title && bytes.Equal(latestVersion.Content, versionToRestore.Content)
		}

		// Restore the selected snapshot to the live document
		document.Title = versionToRestore.Title
		document.Content = versionToRestore.Content
		document.UpdatedBy = userID
		document.CollabVersion++

		if err := s.model.UpdateDocTx(tx, document); err != nil {
			return err
		}

		// If the content is already identical, no new history
		// snapshot is necessary
		if sameAsLatest {
			return nil
		}
		newVersion := entities.Version{
			ID:         uuid.New().String(),
			DocumentID: documentID,
			Version:    nextVersion,
			Title:      versionToRestore.Title,
			Content:    versionToRestore.Content,
			ModifiedBy: userID,
		}

		return s.versionModel.CreateVersion(tx, &newVersion)
	})
}

func (s *documentService) SaveVersion(userID, documentID string) error {
	if _, err := s.CheckEditPermission(userID, documentID); err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		document, err := s.model.GetDocumentForUpdate(tx, documentID)
		if err != nil {
			return err
		}
		latestVersion, err := s.versionModel.GetLatestVersion(tx, documentID)
		if err != nil {
			return err
		}
		// First snapshot for this document.
		nextVersion := 1

		if latestVersion != nil {
			sameTitle := latestVersion.Title == document.Title
			sameContent := bytes.Equal(latestVersion.Content, document.Content)

			if sameTitle && sameContent {
				return nil
			}

			nextVersion = latestVersion.Version + 1
		}
		docVersion := entities.Version{
			ID:         uuid.New().String(),
			DocumentID: document.DocumentID,
			Version:    nextVersion,
			Title:      document.Title,
			Content:    document.Content,
			ModifiedBy: userID,
		}
		return s.versionModel.CreateVersion(tx, &docVersion)
	})
}
