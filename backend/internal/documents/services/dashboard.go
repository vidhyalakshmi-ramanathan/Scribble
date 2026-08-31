package services

import (
	auth "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/models"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/models"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/dto"
	"gorm.io/gorm"
)

type DashboardService interface {
	GetDashboard(userID string) (*dto.DashboardResponse, error)
}

type dashboardService struct {
	db        *gorm.DB
	docModel  models.DocumentModel
	roleModel models.RoleModel
	userModel auth.UserModel
}

func NewDashboardService(db *gorm.DB, docModel models.DocumentModel, roleModel models.RoleModel, userModel auth.UserModel) DashboardService {
	return &dashboardService{
		db:        db,
		docModel:  docModel,
		roleModel: roleModel,
		userModel: userModel,
	}
}

func (s *dashboardService) GetDashboard(userID string) (*dto.DashboardResponse, error) {
	user, err := s.userModel.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	myDocuments, err := s.docModel.GetMyDocs(userID)
	if err != nil {
		return nil, err
	}

	recentlyEdited, err := s.docModel.GetRecentlyEdited(userID, 5)
	if err != nil {
		return nil, err
	}

	sharedDocuments, err := s.docModel.GetSharedDocuments(userID)
	if err != nil {
		return nil, err
	}

	pendingInvites, err := s.roleModel.FindPendingByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	return &dto.DashboardResponse{
		MyDocuments:     myDocuments,
		RecentDocuments: recentlyEdited,
		SharedDocuments: sharedDocuments,
		PendingInvites:  pendingInvites,
	}, nil
}
