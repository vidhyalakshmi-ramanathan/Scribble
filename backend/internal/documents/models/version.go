package models

import (
	"errors"
	"fmt"

	doc "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"gorm.io/gorm"
)

type VersionModel interface {
	CreateVersion(tx *gorm.DB, version *doc.Version) error
	GetVersions(documentID string) (documentVersions []doc.Version, err error)
	GetVersion(documentID string, versionID int) (docVersion *doc.Version, err error)
	GetLatestVersion(documentID string) (version int, err error)
	UpdateContentWithVersion(tx *gorm.DB, documentID string, document *doc.Documents, version int) (bool, error)
}

type versionModel struct {
	db *gorm.DB
}

func NewVersionModel(db *gorm.DB) VersionModel {
	return &versionModel{db: db}
}

func (u *versionModel) CreateVersion(tx *gorm.DB, version *doc.Version) error {
	err := tx.Create(version).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *versionModel) GetVersions(documentID string) (documentVersions []doc.Version, err error) {
	err = u.db.Where("document_id= ?", documentID).Find(&documentVersions).Error
	if err != nil {
		return nil, fmt.Errorf("Error fetching document's versions %w", err)
	}
	return documentVersions, nil
}

func (u *versionModel) GetVersion(documentID string, versionID int) (docVersion *doc.Version, err error) {
	err = u.db.Where("document_id= ? AND version_id= ?", documentID, versionID).First(&docVersion).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Version not found %w", response.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching document by VersionID %w", err)
	}
	return docVersion, nil
}

func (u *versionModel) GetLatestVersion(documentID string) (int, error) {
	var version doc.Version
	err := u.db.Where("document_id= ?", documentID).Order("version_id DESC").First(&version).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("Error fetching version %w", err)
	}
	return version.VersionID, nil
}

func (u *versionModel) UpdateContentWithVersion(tx *gorm.DB, documentID string, document *doc.Documents, version int) (bool, error) {
	result := tx.Model(&doc.Documents{}).
		Where("document_id= ? AND version= ?", documentID, version).
		Updates(map[string]interface{}{
			"content": document.Content,
			"version": gorm.Expr("version + 1"),
		})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}
