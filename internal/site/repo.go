package site

import (
	"context"
	"fmt"
	"time"

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

// get intro from repository
func (r *repo) getIntro(ctx context.Context) (string, error) {
	var data Intro
	err := r.db.
		WithContext(ctx).
		Order("created_at DESC").
		First(&data).
		Error
	if err != nil {
		return "", fmt.Errorf("repo get intro failed: %w", err)
	}

	return data.Content, nil
}

func (r *repo) listAnnouncements(ctx context.Context) ([]Announcement, error) {
	var data []Announcement
	now := time.Now()
	err := r.db.
		WithContext(ctx).
		Model(&Announcement{}).
		Where("online_at <= ? AND offline_at >= ?", now, now).
		Find(&data).
		Error
	if err != nil {
		return nil, fmt.Errorf("repo get announcements failed: %w", err)
	}

	return data, nil
}

func (r *repo) listAnnouncementAdmin(ctx context.Context) ([]Announcement, error) {
	var data []Announcement
	err := r.db.
		WithContext(ctx).
		Model(&Announcement{}).
		Find(&data).
		Error
	if err != nil {
		return nil, fmt.Errorf("repo list announcement failed: %w", err)
	}

	return data, nil
}

type danmakuQuery struct {
	status *DanmakuStatus
}

func (r *repo) listDanmaku(ctx context.Context, query danmakuQuery) ([]Danmaku, error) {
	var danmaku []Danmaku
	db := r.db.WithContext(ctx)

	// 查询范围限制
	if query.status != nil {
		db = db.Where("status = ?", *query.status).Limit(100)
	}

	err := db.
		Order("created_at DESC").
		Find(&danmaku).
		Error
	if err != nil {
		return nil, fmt.Errorf("repo list danmaku failed: %w", err)
	}

	return danmaku, nil
}

func (r *repo) updateDanmakuStatus(ctx context.Context, id int, status DanmakuStatus) error {
	err := r.db.
		WithContext(ctx).
		Model(&Danmaku{}).
		Where("id = ?", id).
		Update("status", status).
		Error
	if err != nil {
		return fmt.Errorf("repo update danmaku status failed: %w", err)
	}

	return nil
}

func (r *repo) createDanmaku(ctx context.Context, d *Danmaku) (int, error) {
	err := r.db.WithContext(ctx).Create(d).Error
	if err != nil {
		return 0, fmt.Errorf("repo create danmaku failed: %w", err)
	}

	return d.ID, nil
}
