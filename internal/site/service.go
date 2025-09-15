package site

import (
	"context"
	"errors"
	"fmt"
	"my_web/backend/internal/zerrors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	db  *repo
	rdb *cache

	cfg func() *Config
}

func NewService(db *gorm.DB, rdb *redis.Client, cfg func() *Config) *Service {
	s := Service{
		cfg: cfg,
	}

	s.db = newRepo(db)
	s.rdb = newCache(rdb, cfg)

	return &s
}

func (s *Service) getIntro(ctx context.Context) (string, error) {
	data, err := s.rdb.getIntro(ctx)
	if err == nil {
		return data, err
	}
	if errors.Is(err, zerrors.CacheMiss) {
	}

	data, err = s.db.getIntro(ctx)
	if err != nil {
		return "", fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	s.rdb.setIntro(ctx, data)
	return data, nil
}

func (s *Service) getAnnouncement(ctx context.Context) ([]Announcement, error) {
	data, err := s.rdb.listAnnouncements(ctx)
	if err == nil {
		return data, err
	}
	if !errors.Is(err, zerrors.CacheMiss) {

	}

	data, err = s.db.listAnnouncements(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	s.rdb.setAnnouncement(ctx, data)
	return data, nil
}

type DanmakuVO struct {
	Content string `json:"content"`
}

type DanmakuAdminVO struct {
	ID        int           `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	SenderID  string        `json:"sender_id"`
	Content   string        `json:"content"`
	Status    DanmakuStatus `json:"status"`
}

func (s *Service) listDanmaku(ctx context.Context) ([]DanmakuVO, error) {
	data, err := s.rdb.getDanmaku(ctx)
	if err == nil {
		return data, err
	}
	if errors.Is(err, zerrors.CacheMiss) {

	}

	status := Approved
	list, err := s.db.listDanmaku(ctx, danmakuQuery{
		status: &status,
	})
	if err != nil {
		return nil, err
	}

	res := make([]DanmakuVO, 0, len(list))
	for _, v := range list {
		res = append(res, DanmakuVO{
			Content: v.Content,
		})
	}

	s.rdb.setDanmaku(ctx, res)
	return res, nil
}

func (s *Service) listDanmakuAdmin(ctx context.Context) ([]DanmakuAdminVO, error) {
	list, err := s.db.listDanmaku(ctx, danmakuQuery{
		status: nil,
	})
	if err != nil {
		return nil, err
	}

	res := make([]DanmakuAdminVO, 0, len(list))
	for _, v := range list {
		res = append(res, DanmakuAdminVO{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			Content:   v.Content,
			SenderID:  v.SenderID,
			Status:    v.Status,
		})
	}
	return res, nil
}

func (s *Service) reviewDanmaku(ctx context.Context, id int, status DanmakuStatus) error {
	if status != Approved && status != Rejected {
		return zerrors.ErrRequest
	}

	err := s.db.updateDanmakuStatus(ctx, id, status)
	if err != nil {
		return fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	return nil
}

func (s *Service) sendDanmaku(ctx context.Context, content, senderID, ip string) (int, error) {
	if content == "" {
		return 0, zerrors.ErrRequest
	}

	d := Danmaku{
		Content:  content,
		SenderID: senderID,
		IP:       ip,
		Status:   Pending,
	}

	id, err := s.db.createDanmaku(ctx, &d)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	return id, nil
}
