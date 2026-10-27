package models

import (
	"errors"
	"fmt"

	doc "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DocumentModel interface {
	CreateDoc(document *doc.Documents) error
	GetDocByID(documentID string) (*doc.Documents, error)
	GetMyAllDocs(userID string) (documents []doc.Documents, err error)
	UpdateDoc(document *doc.Documents) error
	DeleteDoc(documentID string) error
	GetRecentlyEdited(userID string, limit int) ([]doc.Documents, error)
	GetSharedDocuments(userID string) ([]doc.Documents, error)
	GetMyDocs(userID string) (documents []doc.Documents, err error)
	GetDocumentForUpdate(tx *gorm.DB, documentID string) (*doc.Documents, error)
	UpdateDocTx(tx *gorm.DB, document *doc.Documents) error
}

type documentModel struct {
	db *gorm.DB
}

func NewDocumentModel(db *gorm.DB) DocumentModel {
	return &documentModel{db: db}
}

func (u *documentModel) CreateDoc(document *doc.Documents) error {
	err := u.db.Create(document).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *documentModel) GetDocByID(documentID string) (*doc.Documents, error) {
	var docDetails doc.Documents
	err := u.db.Where("document_id= ?", documentID).First(&docDetails).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Document not found %w", response.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching document details by ID %w", err)
	}
	return &docDetails, nil
}

func (u *documentModel) GetMyAllDocs(userID string) ([]doc.Documents, error) {
	var documents []doc.Documents
	err := u.db.
		Where(
			"owner_id= ? OR document_id IN (?)", userID,
			u.db.Table("permissions").Select("document_id").Where("user_id= ?", userID)).
		Order("updated_at DESC").Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("Error fetching documents details %w", err)
	}
	return documents, nil
}

func (u *documentModel) UpdateDoc(document *doc.Documents) error {
	return u.db.Model(&doc.Documents{}).Where("document_id= ?", document.DocumentID).
		Updates(map[string]interface{}{
			"title":      document.Title,
			"content":    document.Content,
			"category":   document.Category,
			"updated_by": document.UpdatedBy,
			"version":    document.Version,
		}).Error
}

func (u *documentModel) DeleteDoc(documentID string) error {
	err := u.db.Where("document_id= ?", documentID).Delete(&doc.Documents{})
	if err.Error != nil {
		return fmt.Errorf("Error delete document %w", err.Error)
	}
	if err.RowsAffected == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (u *documentModel) GetMyDocs(userID string) (documents []doc.Documents, err error) {
	err = u.db.Where("owner_id= ?)", userID).Order("updated_at DESC").Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("Error fetching documents details by OwnerID %w", err)
	}
	return documents, nil

}

func (u *documentModel) GetRecentlyEdited(userID string, limit int) ([]doc.Documents, error) {
	var documents []doc.Documents

	err := u.db.Joins("LEFT JOIN permissions p ON p.document_id = documents.document_id AND p.user_id = ?", userID).
		Where("documents.owner_id = ? OR p.user_id IS NOT NULL", userID).
		Order("documents.updated_at DESC").
		Limit(limit).
		Find(&documents).Error

	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (u *documentModel) GetSharedDocuments(userID string) ([]doc.Documents, error) {
	var documents []doc.Documents

	err := u.db.Joins("JOIN permissions p ON p.document_id = documents.document_id").
		Where("p.user_id = ? AND documents.owner_id != ?", userID, userID).
		Order("documents.updated_at DESC").
		Find(&documents).Error

	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (u *documentModel) GetDocumentForUpdate(tx *gorm.DB, documentID string) (*doc.Documents, error) {
	var document doc.Documents

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("document_id = ?", documentID).
		First(&document).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Document not found %w", response.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching document details by ID %w", err)
	}
	return &document, nil
}

func (u *documentModel) UpdateDocTx(tx *gorm.DB, document *doc.Documents) error {
	result := tx.Model(&doc.Documents{}).
		Where("document_id = ?", document.DocumentID).
		Updates(map[string]interface{}{
			"title":      document.Title,
			"content":    document.Content,
			"updated_by": document.UpdatedBy,
			"version":    document.Version,
		})
	if result.Error != nil {
		return fmt.Errorf("error updating document: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return response.ErrNotFound
	}
	return nil
}
