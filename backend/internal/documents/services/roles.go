package services

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	invite "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/entities"
	auth "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/models"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/models"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/enums"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"gorm.io/gorm"
)

type RoleService interface {
	GetMemberDetails(userID, targetID, documentID string) (*entities.MemberDetailRes, error)
	GetMembersList(userID, documentID string) ([]entities.MemberDetailRes, error)
	UpdateRole(userID, memberID, documentID string, role enums.Roles) error
	DeleteMember(userID, memberID, documentID string) error
	CheckAdminorOwner(userID, docID string) (*entities.Documents, error)
	CheckMember(userID, docID string) (*entities.Documents, error)
	CreateInvite(userID, documentID string, invitedEmail string, role enums.Roles) error
	GetMyInvites(userID string) (invites []invite.DocumentInvites, err error)
	GetSentInvites(userID, docID string) (invites []invite.DocumentInvites, err error)
	AcceptInvite(userID, inviteID string) error
}

type roleService struct {
	rolemodel models.RoleModel
	userModel auth.UserModel
	docModel  models.DocumentModel
	db        *gorm.DB
}

func NewRoleService(rolemodel models.RoleModel, userModel auth.UserModel, docModel models.DocumentModel, db *gorm.DB) RoleService {
	return &roleService{
		rolemodel: rolemodel,
		userModel: userModel,
		docModel:  docModel,
		db:        db,
	}
}

func (s *roleService) GetMemberDetails(userID, targetID, documentID string) (*entities.MemberDetailRes, error) {
	document, err := s.CheckMember(userID, documentID)
	if errors.Is(err, response.ErrUserNotMember) {
		return nil, response.ErrForbidden
	}
	if err != nil {
		return nil, err
	}
	if document.OwnerID == targetID {
		return &entities.MemberDetailRes{
			UserID:     targetID,
			DocumentID: document.DocumentID,
			Role:       enums.Owner,
		}, nil
	}
	var memberDetails *entities.Permissions
	memberDetails, err = s.rolemodel.GetMemberDetails(targetID, documentID)
	if errors.Is(err, response.ErrNotFound) {
		return nil, response.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	memberDetailsRes := &entities.MemberDetailRes{
		UserID:     memberDetails.UserID,
		DocumentID: memberDetails.DocumentID,
		Role:       memberDetails.Role,
	}
	return memberDetailsRes, nil
}

func (s *roleService) GetMembersList(userID, documentID string) ([]entities.MemberDetailRes, error) {
	_, err := s.CheckMember(userID, documentID)
	if errors.Is(err, response.ErrUserNotMember) {
		return nil, response.ErrForbidden
	}
	if err != nil {
		return nil, err
	}

	var members []entities.Permissions
	members, err = s.rolemodel.GetMembersList(documentID)
	if err != nil {
		return nil, err
	}
	var membersDetailsRes []entities.MemberDetailRes
	for _, member := range members {
		res := entities.MemberDetailRes{
			UserID:     member.UserID,
			DocumentID: member.DocumentID,
			Role:       member.Role,
		}
		membersDetailsRes = append(membersDetailsRes, res)
	}
	return membersDetailsRes, nil
}

func (s *roleService) UpdateRole(userID, memberID, documentID string, role enums.Roles) error {
	document, err := s.CheckAdminorOwner(userID, documentID)
	if err != nil {
		return err
	}
	// owner role cannot be changed
	if document.OwnerID == memberID {
		return response.ErrForbidden
	}
	if err := s.ValidationRole(role); err != nil {
		return err
	}
	member, err := s.rolemodel.GetMemberDetails(memberID, documentID)
	if err != nil {
		return err
	}
	//admin cannot change another admin
	if document.OwnerID != userID && member.Role != enums.Admin {
		return response.ErrForbidden
	}
	return s.rolemodel.UpdateRole(memberID, documentID, role)
}

func (s *roleService) DeleteMember(userID, memberID, documentID string) error {
	document, err := s.CheckAdminorOwner(userID, documentID)
	if err != nil {
		return err
	}
	// owner role cannot be deleted
	if document.OwnerID == memberID {
		return response.ErrForbidden
	}
	member, err := s.rolemodel.GetMemberDetails(memberID, documentID)
	if err != nil {
		return err
	}
	//admin cannot delete another admin
	if document.OwnerID != userID && member.Role != enums.Admin {
		return response.ErrForbidden
	}
	return s.rolemodel.DeleteMember(memberID, documentID)
}

// CheckUserPermission checks who is making the request admin/owner
func (s *roleService) CheckAdminorOwner(userID, docID string) (*entities.Documents, error) {
	document, err := s.docModel.GetDocByID(docID)
	if err != nil {
		return nil, err
	}
	if err == response.ErrNotFound {
		return nil, response.ErrDocumentNotExists
	}
	if document.OwnerID == userID {
		return document, nil
	}

	user, err := s.rolemodel.GetMemberDetails(userID, docID)
	if err != nil {
		return nil, err
	}
	if user.Role != enums.Admin {
		return nil, response.ErrForbidden
	}
	return document, nil
}

func (s *roleService) ValidationRole(role enums.Roles) error {
	if role != enums.Admin && role != enums.Editor && role != enums.Viewer {
		return response.ErrInvalidRequest
	}
	return nil
}

// CheckMember checks who is making the request
func (s *roleService) CheckMember(userID, docID string) (*entities.Documents, error) {
	document, err := s.docModel.GetDocByID(docID)
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

func (s *roleService) CreateInvite(userID, documentID string, invitedEmail string, role enums.Roles) error {
	// Owner/Admin and does document exist?
	_, err := s.CheckAdminorOwner(userID, documentID)
	if err != nil {
		return err
	}
	if invitedEmail == "" {
		return &response.ValidationError{
			Message: "Email cannot be empty",
			Code:    http.StatusBadRequest,
		}
	}
	if err := s.ValidationRole(role); err != nil {
		return err
	}
	invitedUser, err := s.userModel.GetUserByEmail(invitedEmail)
	if err != nil {
		return err
	}
	if invitedUser.UserID == userID {
		return response.ErrCannotInviteSelf
	}

	// is the invited user already a member
	member, err := s.rolemodel.GetMemberDetails(invitedUser.UserID, documentID)
	if err == nil && member != nil {
		return response.ErrUserAlreadyMember
	}
	if err != nil && !errors.Is(err, response.ErrNotFound) {
		return err
	}

	invites, err := s.rolemodel.FindPendingByEmail(invitedEmail)
	if err != nil {
		return err
	}
	for _, invite := range invites {
		if invite.DocumentID == documentID {
			return response.ErrInviteAlreadyExists
		}
	}

	invite := &invite.DocumentInvites{
		ID:           uuid.New().String(),
		DocumentID:   documentID,
		InvitedEmail: invitedEmail,
		Role:         role,
		InvitedBy:    userID,
		Status:       enums.InvitePending,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}
	return s.rolemodel.AddInvite(invite)
}

func (s *roleService) GetMyInvites(userID string) (invites []invite.DocumentInvites, err error) {
	invitedUser, err := s.userModel.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return s.rolemodel.FindPendingByEmail(invitedUser.Email)
}

func (s *roleService) GetSentInvites(userID, docID string) (invites []invite.DocumentInvites, err error) {
	_, err = s.CheckAdminorOwner(userID, docID)
	if err != nil {
		return nil, err
	}
	return s.rolemodel.FindInvitesByDocumentID(docID)
}

func (s *roleService) AcceptInvite(userID, inviteID string) error {
	invite, err := s.rolemodel.FindInviteByID(inviteID)
	if errors.Is(err, response.ErrNotFound) {
		return response.ErrInviteNotExists
	}
	if err != nil {
		return err
	}

	if invite.Status != enums.InvitePending {
		return response.ErrInviteNotPendingOrExpired
	}
	if time.Now().After(invite.ExpiresAt) {
		return response.ErrInviteNotPendingOrExpired
	}

	user, err := s.userModel.GetUserByID(userID)
	if err != nil {
		return err
	}
	if !strings.EqualFold(user.Email, invite.InvitedEmail) {
		return response.ErrForbidden
	}

	// is the invited user already a member
	member, err := s.rolemodel.GetMemberDetails(userID, invite.DocumentID)
	if err == nil && member != nil {
		return response.ErrUserAlreadyMember
	}
	if err != nil && !errors.Is(err, response.ErrNotFound) {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		member := &entities.Permissions{
			ID:         uuid.New().String(),
			DocumentID: invite.DocumentID,
			UserID:     userID,
			Role:       invite.Role,
		}
		if err := s.rolemodel.AddMemberRole(tx, member); err != nil {
			return err
		}

		if err = s.rolemodel.UpdateInviteStatus(tx, invite.ID, enums.InviteAccepted); err != nil {
			return err
		}
		return nil
	})
}
