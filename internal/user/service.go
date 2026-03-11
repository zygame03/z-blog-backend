package user

import (
	"context"
	"my_web/backend/internal/logger"
	"my_web/backend/internal/middleware"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db   *repo
	rdb  *cache
	conf func() *Config
}

func NewService(db *gorm.DB, rdb *redis.Client, conf func() *Config) *Service {
	return &Service{
		db:   newRepo(db),
		rdb:  newCache(rdb),
		conf: conf,
	}
}

func (s *Service) getProfile(ctx context.Context, id int) (*Profile, error) {
	data, err := s.rdb.getProfile(ctx, id)

	if err == nil {
		logger.Info(
			"get profile from cache",
			zap.Int("id", id),
		)
		return data, nil
	}
	if err != ErrCacheMiss {
		logger.Warn(
			"cache get profile failed",
			zap.Error(err),
		)
	}

	data, err = s.db.getProfile(id)
	if err != nil {
		logger.Error(
			"repo get profile failed",
			zap.Int("id", id),
			zap.Error(err),
		)
		return nil, err
	}

	err = s.rdb.setProfile(ctx, id, data)
	if err != nil {
		logger.Warn(
			"cache set profile failed",
			zap.Int("id", id),
			zap.Error(err),
		)
	}
	return data, nil
}

func (s *Service) Register(ctx context.Context, username, password string) error {
	user, err := s.db.getUserByUsername(username)
	if err != nil {
		return err
	}
	if user != nil {
		return ErrUserAlreadyExist
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := &User{
		Username:   username,
		Password:   string(hashedPassword),
		Permission: 0,
	}
	return s.db.createUser(newUser)
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	// 用户名查找
	user, err := s.db.getUserByUsername(username)
	if err != nil {
		return "", err
	}
	// 未找到用户
	if user == nil {
		return "", ErrUserNotFound
	}

	// 找到用户，对比密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}
	ttl := 24 * time.Hour
	if s.conf != nil {
		conf := s.conf()
		if conf != nil && conf.TokenTTL > 0 {
			ttl = conf.TokenTTL
		}
	}

	// 获取token
	token, err := middleware.GenerateToken(user.ID, user.Username, ttl)
	if err != nil {
		return "", err
	}
	return token, nil
}
