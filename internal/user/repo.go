package user

import (
	"errors"
	"my_web/backend/internal/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func newRepo(db *gorm.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) getProfile(id int) (*Profile, error) {
	var profile Profile
	result := r.db.
		Model(&Profile{}).
		Where("id = ?", id).
		First(&profile)
	if result.Error != nil {
		logger.Error(
			"db get profile failed",
			zap.Int("id", id),
		)
		return nil, result.Error
	}

	return &profile, nil
}

func (r *repo) getUserByUsername(username string) (*User, error) {
	var user User
	result := r.db.
		Model(&User{}).
		Where("username = ?", username).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Error(
			"db get user failed",
			zap.String("username", username),
			zap.Error(result.Error),
		)
		return nil, result.Error
	}
	return &user, nil
}

func (r *repo) createUser(user *User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		logger.Error(
			"db create user failed",
			zap.String("username", user.Username),
			zap.Error(result.Error),
		)
		return result.Error
	}
	return nil
}
