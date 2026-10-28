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
	GetVersions(documentID string) (documentVersions []doc.VersionsList, err error)
	GetVersion(documentID string, versionID int) (*doc.Version, error)
	GetLatestVersion(tx *gorm.DB, documentID string) (*doc.Version, error)
	UpdateContentByVersion(documentID string, document *doc.Documents, version int) (bool, error)
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
		return fmt.Errorf("error creating document version: %w", err)
	}
	return nil
}

func (u *versionModel) GetVersions(documentID string) (documentVersions []doc.VersionsList, err error) {
	err = u.db.Raw(`
        SELECT
            v.id AS id,
            v.document_id AS document_id,
            v.version_id AS version_id,
            v.title AS title,
            v.modified_by AS modified_by_user_id,
            COALESCE(u.user_name, 'Unknown user') AS modified_by_user_name,
            v.created_at AS created_at
        FROM versions v
        LEFT JOIN users u ON u.user_id = v.modified_by
        WHERE v.document_id = ?
        ORDER BY v.created_at DESC
    `, documentID).Scan(&documentVersions).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching document versions: %w", err)
	}
	return documentVersions, nil
}

func (u *versionModel) GetVersion(documentID string, versionID int) (*doc.Version, error) {
	var docVersion doc.Version
	err := u.db.Where("document_id= ? AND version_id= ?", documentID, versionID).First(&docVersion).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Version not found %w", response.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching document by VersionID %w", err)
	}
	return &docVersion, nil
}

func (u *versionModel) GetLatestVersion(tx *gorm.DB, documentID string) (*doc.Version, error) {
	var latestVersion doc.Version
	err := tx.Where("document_id= ?", documentID).Order("version_id DESC").First(&latestVersion).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching latest version %w", err)
	}
	return &latestVersion, nil
}

func (u *versionModel) UpdateContentByVersion(documentID string, document *doc.Documents, version int) (bool, error) {
	result := u.db.Model(&doc.Documents{}).
		Where("document_id= ? AND version= ?", documentID, version).
		Updates(map[string]interface{}{
			"content":    document.Content,
			"version":    gorm.Expr("version + 1"),
			"updated_by": document.UpdatedBy,
		})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}
