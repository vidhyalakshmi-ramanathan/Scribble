package models

import (
	"errors"
	"fmt"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	doc "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"gorm.io/gorm"
)

type VersionModel interface {
	CreateVersion(tx *gorm.DB, version *doc.Version) error
	GetVersions(documentID string) ([]doc.VersionsList, error)
	GetVersion(documentID string, versionID int) (*doc.VersionDetailRes, error)
	GetLatestVersion(tx *gorm.DB, documentID string) (*doc.Version, error)
	ApplyCollaborativeEdit(documentID string, document *doc.Documents, baseCollabVersion int) (bool, error)
	UpdateSnapshotVersionTx(tx *gorm.DB, documentID string, version int) error
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

func (u *versionModel) GetVersions(documentID string) ([]doc.VersionsList, error) {
	documentVersions := make([]doc.VersionsList, 0)
	err := u.db.Raw(`
        SELECT
            v.id AS id,
            v.document_id AS docID,
            v.version AS version,
            v.title AS title,
			v.schema_ver AS schemaVer,
            v.modified_by AS modifiedByUserID,
            COALESCE(u.user_name, 'Unknown user') AS modifiedByUserName,
            v.created_at AS createdAt
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

func (u *versionModel) GetVersion(documentID string, version int) (*doc.VersionDetailRes, error) {
	var docVersion doc.VersionDetailRes
	err := u.db.Raw(`
        SELECT
            v.id,
            v.document_id AS docID,
            v.version,
            v.title,
            v.content,
            v.schema_ver AS schemaVer,
            v.modified_by AS modifiedByUserID,
            COALESCE(u.user_name, 'Unknown user') AS modifiedByUserName,
            v.created_at AS createdAt
        FROM versions v
        LEFT JOIN users u ON u.user_id = v.modified_by
        WHERE v.document_id = ? AND v.version = ?
    `, documentID, version).Scan(&docVersion).Error

	if err != nil {
		return nil, fmt.Errorf("Error fetching document by version %w", err)
	}
	if docVersion.ID == "" {
		return nil, fmt.Errorf("version not found: %w", response.ErrNotFound)
	}
	return &docVersion, nil
}

func (u *versionModel) GetLatestVersion(tx *gorm.DB, documentID string) (*doc.Version, error) {
	var latestVersion doc.Version
	err := tx.Where("document_id= ?", documentID).Order("version DESC").First(&latestVersion).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching latest document version %w", err)
	}
	return &latestVersion, nil
}

func (u *versionModel) ApplyCollaborativeEdit(documentID string, document *doc.Documents, baseCollabVersion int) (bool, error) {
	result := u.db.Model(&doc.Documents{}).
		Where("document_id= ? AND collab_version= ?", documentID, baseCollabVersion).
		Updates(map[string]interface{}{
			"content":        document.Content,
			"collab_version": gorm.Expr("collab_version + 1"),
			"updated_by":     document.UpdatedBy,
		})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}

func (u *versionModel) UpdateSnapshotVersionTx(tx *gorm.DB, documentID string, version int) error {
	result := tx.Model(&entities.Documents{}).
		Where("document_id = ?", documentID).
		Update("version", version)

	if result.Error != nil {
		return fmt.Errorf("Error Update document snapshot version: %w", result.Error)
	}

	if result.RowsAffected != 1 {
		return response.ErrNotFound
	}

	return nil
}
