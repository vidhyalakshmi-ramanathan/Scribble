package models

import (
	"errors"
	"fmt"
	"time"

	auth "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/dto"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/enums"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"gorm.io/gorm"
)

type RoleModel interface {
	AddMemberRole(tx *gorm.DB, user *entities.Permissions) error
	GetMemberDetails(userID, documentID string) (*entities.Permissions, error)
	GetMembersList(documentID string) ([]entities.Permissions, error)
	UpdateRole(userID, documentID string, role enums.Roles) error
	DeleteMember(userID, documentID string) error
	AddInvite(invite *auth.DocumentInvites) error
	FindInviteByID(id string) (*auth.DocumentInvites, error)
	FindInvitesByDocumentID(documentID string) ([]auth.DocumentInvites, error)
	FindPendingByEmail(email string) ([]auth.DocumentInvites, error)
	UpdateInviteStatus(tx *gorm.DB, id string, status enums.Invites) error
	GetPendingInvites(email string) ([]dto.PendingInvites, error)
}

type roleModel struct {
	db *gorm.DB
}

func NewRoleModel(db *gorm.DB) RoleModel {
	return &roleModel{db: db}
}

func (u *roleModel) AddMemberRole(tx *gorm.DB, user *entities.Permissions) error {
	err := tx.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

// GetMemberDetails function gives the specific user'role of a document
func (u *roleModel) GetMemberDetails(userID, documentID string) (*entities.Permissions, error) {
	var user *entities.Permissions
	err := u.db.Where("user_id= ? AND document_id= ?", userID, documentID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Document not found %w", response.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching user details %w", err)
	}
	return user, nil
}

func (u *roleModel) GetMembersList(documentID string) ([]entities.Permissions, error) {
	var user []entities.Permissions
	err := u.db.Where("document_id= ?", documentID).Find(&user).Error
	if err != nil {
		return nil, fmt.Errorf("Error fetching users %w", err)
	}
	return user, nil
}

func (u *roleModel) UpdateRole(userID, documentID string, role enums.Roles) error {
	err := u.db.Model(&entities.Permissions{}).Where("user_id= ? AND document_id= ?", userID, documentID).Updates(map[string]interface{}{
		"role": role,
	}).Error
	if err != nil {
		return fmt.Errorf("Error updating role %w", err)
	}
	return nil
}

func (u *roleModel) DeleteMember(userID, documentID string) error {
	err := u.db.Where("user_id= ? AND document_id= ?", userID, documentID).Delete(&entities.Permissions{})
	if err.Error != nil {
		return fmt.Errorf("Error delete member %w", err.Error)
	}
	if err.RowsAffected == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (u *roleModel) AddInvite(invite *auth.DocumentInvites) error {
	err := u.db.Create(invite).Error
	if err != nil {
		return fmt.Errorf("Error adding invite %w", err)
	}
	return nil
}

func (u *roleModel) FindInviteByID(id string) (*auth.DocumentInvites, error) {
	var invite *auth.DocumentInvites
	err := u.db.Where("id= ?", id).First(&invite).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Invite not found %w", response.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching user invite %w", err)
	}
	return invite, nil
}

func (u *roleModel) FindInvitesByDocumentID(documentID string) ([]auth.DocumentInvites, error) {
	var invites []auth.DocumentInvites
	err := u.db.Where("document_id= ?", documentID).Order("created_at DESC").Find(&invites).Error
	if err != nil {
		return nil, fmt.Errorf("Error fetching users %w", err)
	}
	return invites, nil
}

func (u *roleModel) FindPendingByEmail(email string) ([]auth.DocumentInvites, error) {
	var invites []auth.DocumentInvites
	err := u.db.Where("email= ? AND status= ?", email, enums.InvitePending).Find(&invites).Error
	if err != nil {
		return nil, fmt.Errorf("Error fetching pending invites %w", err)
	}
	return invites, nil
}

func (u *roleModel) UpdateInviteStatus(tx *gorm.DB, id string, status enums.Invites) error {
	err := tx.Model(&auth.DocumentInvites{}).Where("id= ?", id).Update("status", status).Error
	if err != nil {
		return fmt.Errorf("Error updating invite status %w", err)
	}
	return nil
}

func (u *roleModel) GetPendingInvites(email string) ([]dto.PendingInvites, error) {
	invites := make([]dto.PendingInvites, 0)

	err := u.db.Model(&auth.DocumentInvites{}).
		Select(`
            i.id AS invitationID,
            i.document_id,
            d.title AS title,
            inviter.user_name AS invitedBy,
            inviter.user_id AS invitedByID,
            i.role,
            i.expires_at AS ,
            i.created_at AS 
        `).
		Table("document_invites i").
		Joins("JOIN documents d ON d.document_id = i.document_id").
		Joins("LEFT JOIN users inviter ON inviter.user_id = i.invited_by").
		Where("LOWER(i.invited_email) = LOWER(?)", email).
		Where("i.status = ?", enums.Invites("pending")).
		Where("i.expires_at > ?", time.Now()).
		Order("i.created_at DESC").
		Scan(&invites).Error

	if err != nil {
		return nil, fmt.Errorf("fetching pending invitations: %w", err)
	}

	return invites, nil
}
