package models

import (
	"errors"
	"fmt"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/entities"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"gorm.io/gorm"
)

type UserModel interface {
	Exists(email string) (exists bool, err error)
	SignUp(user *entities.User) error
	GetUserByEmail(email string) (userDetails *entities.User, err error)
	GetUserByID(userID string) (userDetails *entities.User, err error)
}

type userModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) UserModel {
	return &userModel{db: db}
}

func (u *userModel) Exists(email string) (exists bool, err error) {
	var count int64
	err = u.db.Model(&entities.User{}).Where("email= ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("Error fetching user details %w", err)
	}
	return count > 0, nil
}

func (u *userModel) SignUp(user *entities.User) error {
	err := u.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *userModel) GetUserByEmail(email string) (userDetails *entities.User, err error) {
	err = u.db.Where("email= ?", email).First(&userDetails).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("User not found %w", response.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching user details by email %w", err)
	}
	return userDetails, nil
}

func (u *userModel) GetUserByID(userID string) (userDetails *entities.User, err error) {
	err = u.db.Where("user_id= ?", userID).First(&userDetails).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("User not found %w", response.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching user details by ID %w", err)
	}
	return userDetails, nil
}
