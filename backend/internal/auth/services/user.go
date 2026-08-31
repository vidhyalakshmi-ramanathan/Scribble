package services

import (
	"os"
	"strconv"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/entities/packets"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/models"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUpUser(user *entities.User) error
	Login(user *packets.LoginReq) (*packets.LoginResponse, error)
}

type userService struct {
	model models.UserModel
}

func NewUserService(model models.UserModel) UserService {
	return &userService{model: model}
}

func (s *userService) SignUpUser(user *entities.User) error {
	if user.UserName == "" || user.Email == "" || user.Password == "" {
		return &response.ValidationError{
			Message: "username or email or password cannot be empty",
		}
	}
	exists, err := s.model.Exists(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return response.ErrUserAlreadyExists
	}

	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), cost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.model.SignUp(user)
}

func (s *userService) Login(user *packets.LoginReq) (*packets.LoginResponse, error) {
	userDetails, err := s.model.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(user.Password))
	if err != nil {
		return nil, response.ErrInvalidPassword
	}

	token, err := jwt.GenerateToken(userDetails.UserID, userDetails.UserName)
	if err != nil {
		return nil, response.ErrUnauthorizedAccess
	}
	loginResponse := packets.LoginResponse{
		User: packets.UserDataResponse{
			UserID:   userDetails.UserID,
			Username: userDetails.UserName,
			Email:    userDetails.Email,
		},
		Token: token,
	}
	return &loginResponse, nil
}
