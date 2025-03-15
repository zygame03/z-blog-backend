package service

import (
	"errors"
	"my_web/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

type Profile struct {
	models.Model
	Avatar    string `json:"avatar"`
	Name      string `json:"name"`
	Signature string `json:"signature"`
}

func (u *UserService) authenticate(username, password string) (int, error) {
	var user models.User

	result := u.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return 0, result.Error
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return user.ID, errors.New("invalid password")
	}

	return user.ID, nil
}

func (u *UserService) Register(username, password string) (int, error) {
	var user models.User

	result := u.DB.Where("username = ?", username).First(&user)
	if result.Error == nil {
		return 0, errors.New("用户名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user = models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	err = user.Create(u.DB)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (u *UserService) Login(username, password string) (int, error) {
	return u.authenticate(username, password)
}

func (u *UserService) GetUser() (models.User, error) {
	var user models.User

	result := u.DB.First(&user, 1)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (u *UserService) GetProfile() (Profile, error) {
	var profile Profile

	result := u.DB.Model(&models.User{}).
		Select("id, created_at, updated_at, avatar, name, signature").
		First(&profile)
	if result.Error != nil {
		return profile, result.Error
	}

	return profile, nil
}
