package models

import (
	"errors"
	"fmt"

	doc "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"gorm.io/gorm"
)

type DocumentModel interface {
	CreateDoc(tx *gorm.DB, document *doc.Documents) error
	GetDocByID(documentID string) (docDetails *doc.Documents, err error)
	GetMyAllDocs(userID string) (documents []doc.Documents, err error)
	UpdateDoc(tx *gorm.DB, document *doc.Documents) error
	DeleteDoc(documentID string) error
	GetRecentlyEdited(userID string, limit int) ([]doc.Documents, error)
	GetSharedDocuments(userID string) ([]doc.Documents, error)
	GetMyDocs(userID string) (documents []doc.Documents, err error)
}

type documentModel struct {
	db *gorm.DB
}

func NewDocumentModel(db *gorm.DB) DocumentModel {
	return &documentModel{db: db}
}

func (u *documentModel) CreateDoc(tx *gorm.DB, document *doc.Documents) error {
	err := tx.Create(document).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *documentModel) GetDocByID(documentID string) (docDetails *doc.Documents, err error) {
	err = u.db.Where("document_id= ?", documentID).First(&docDetails).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Document not found %w", response.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching document details by ID %w", err)
	}
	return docDetails, nil
}

func (u *documentModel) GetMyAllDocs(userID string) (documents []doc.Documents, err error) {
	err = u.db.Where("owner_id= ? OR document_id IN (SELECT document_id FROM permissions WHERE user_id= ?)", userID).Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("Error fetching documents details %w", err)
	}
	return documents, nil

}

func (u *documentModel) UpdateDoc(tx *gorm.DB, document *doc.Documents) error {
	return tx.Model(&doc.Documents{}).Where("document_id= ?", document.DocumentID).
		Updates(map[string]interface{}{
			"title":      document.Title,
			"content":    document.Content,
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
