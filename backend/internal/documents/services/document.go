package services

import (
	"github.com/google/uuid"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/models"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/enums"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"gorm.io/gorm"
)

type DocumentService interface {
	CreateDoc(userID string, document *entities.Documents) error
	GetDocByID(userID, documentID string) (docDetails *entities.Documents, err error)
	GetMyAllDocs(userID string) (documents []entities.Documents, err error)
	UpdateDoc(userID string, document *entities.Documents) error
	DeleteDoc(userID, documentID string) error
	CheckEditPermission(userID, docID string) (*entities.Documents, error)
	UpdateContentWithVersion(userID string, document *entities.Documents, baseVersion int) error
	GetVersion(userID, documentID string, version int) (versionDetail *entities.Version, err error)
	GetVersions(userID, documentID string) (versions []entities.Version, err error)
	RestoreVersion(userID, documentID string, version int) error
}

type documentService struct {
	model        models.DocumentModel
	rolemodel    models.RoleModel
	versionModel models.VersionModel
	db           *gorm.DB
}

func NewDocumentService(model models.DocumentModel, rolemodel models.RoleModel, versionModel models.VersionModel, db *gorm.DB) DocumentService {
	return &documentService{
		model:        model,
		rolemodel:    rolemodel,
		versionModel: versionModel,
		db:           db,
	}
}

func (s *documentService) CreateDoc(userID string, document *entities.Documents) error {
	if document.Title == "" {
		return &response.ValidationError{
			Message: "title cannot be empty",
		}
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.model.CreateDoc(tx, document); err != nil {
			return err
		}

		docVersion := entities.Version{
			ID:         uuid.New().String(),
			DocumentID: document.DocumentID,
			VersionID:  1,
			Content:    document.Content,
			ModifiedBy: userID,
		}
		if err := s.versionModel.CreateVersion(tx, &docVersion); err != nil {
			return err
		}
		return nil
	})
}

func (s *documentService) GetDocByID(userID, documentID string) (docDetails *entities.Documents, err error) {
	docDetails, err = s.CheckAccess(userID, documentID)
	if err != nil {
		return nil, err
	}
	return docDetails, nil
}

func (s *documentService) GetMyAllDocs(userID string) (documents []entities.Documents, err error) {
	documents, err = s.model.GetMyAllDocs(userID)
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (s *documentService) UpdateDoc(userID string, document *entities.Documents) error {
	if document.Title == "" {
		return &response.ValidationError{
			Message: "title cannot be empty",
		}
	}
	oldDocument, err := s.CheckEditPermission(userID, document.DocumentID)
	if err != nil {
		return err
	}

	nextVersion := oldDocument.Version + 1

	oldDocument.Title = document.Title
	oldDocument.Content = document.Content
	oldDocument.UpdatedBy = userID
	oldDocument.Version = nextVersion

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.model.UpdateDoc(tx, oldDocument); err != nil {
			return err
		}

		docVersion := entities.Version{
			ID:         uuid.New().String(),
			DocumentID: oldDocument.DocumentID,
			VersionID:  nextVersion,
			Content:    oldDocument.Content,
			ModifiedBy: userID,
		}
		if err := s.versionModel.CreateVersion(tx, &docVersion); err != nil {
			return err
		}
		return nil
	})
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
		return nil, err
	}
	if err == response.ErrNotFound {
		return nil, response.ErrDocumentNotExists
	}
	return document, nil
}

func (s *documentService) CheckAccess(userID, docID string) (*entities.Documents, error) {
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

func (s *documentService) UpdateContentWithVersion(userID string, document *entities.Documents, baseVersion int) error {

	return s.db.Transaction(func(tx *gorm.DB) error {
		updated, err := s.versionModel.UpdateContentWithVersion(tx, document.DocumentID, document, baseVersion)
		if err != nil {
			return err
		}
		if !updated {
			return response.ErrVersionConflict
		}
		nextVersion := baseVersion + 1

		docVersion := entities.Version{
			ID:         uuid.New().String(),
			DocumentID: document.DocumentID,
			VersionID:  nextVersion,
			Content:    document.Content,
			ModifiedBy: userID,
		}
		if err := s.versionModel.CreateVersion(tx, &docVersion); err != nil {
			return err
		}
		return nil
	})
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

func (s *documentService) GetVersions(userID, documentID string) (versions []entities.Version, err error) {
	_, err = s.CheckAccess(userID, documentID)
	if err != nil {
		return nil, err
	}
	versions, err = s.versionModel.GetVersions(documentID)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func (s *documentService) RestoreVersion(userID, documentID string, version int) error {
	document, err := s.CheckEditPermission(userID, documentID)
	if err != nil {
		return err
	}

	oldVersion, err := s.versionModel.GetVersion(documentID, version)
	if err != nil {
		return err
	}

	newVersion := document.Version + 1

	return s.db.Transaction(func(tx *gorm.DB) error {
		document.Content = oldVersion.Content
		document.Version = newVersion
		document.UpdatedBy = userID

		if err := s.model.UpdateDoc(tx, document); err != nil {
			return err
		}

		docVersion := entities.Version{
			ID:         uuid.New().String(),
			DocumentID: documentID,
			VersionID:  newVersion,
			Content:    oldVersion.Content,
			ModifiedBy: userID,
		}

		if err := s.versionModel.CreateVersion(tx, &docVersion); err != nil {
			return err
		}

		return nil
	})
}
