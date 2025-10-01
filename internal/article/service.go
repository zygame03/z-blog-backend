package article

import (
	"context"
	"errors"
	"fmt"
	"my_web/backend/internal/logger"
	"my_web/backend/internal/zerrors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ArticleService struct {
	db  *repo
	rdb *cache

	cfg func() *Config
}

func NewService(db *gorm.DB, rdb *redis.Client, cfg func() *Config) *ArticleService {
	service := &ArticleService{
		cfg: cfg,
	}

	service.db = newRepo(db)
	service.rdb = newCache(rdb, cfg)

	return service
}

func (s *ArticleService) RegisterCron(cron *cron.Cron) {
	_, err := cron.AddFunc("@every "+s.cfg().SyncInterval.String(), s.syncArticleViews)
	if err != nil {
		return
	}

	logger.Info(
		"cron job registered",
		zap.String("job", "syncArticleViews"),
		zap.Duration("interval", s.cfg().SyncInterval),
	)
}

func (s *ArticleService) syncArticleViews() {
	ctx := context.Background()
	ids, err := s.db.listIDs(ctx)
	if err != nil {
		logger.Error(
			"load article ids for view sync failed",
			zap.Error(err),
		)
		return
	}

	var (
		total   = len(ids)
		success int
		failed  int
	)
	for _, id := range ids {
		num, err := s.rdb.getViewUV(ctx, id)
		if err != nil {
			failed++
			continue
		}
		if num == 0 {
			success++
			continue
		}
		err = s.db.incrementViews(ctx, id, num)
		if err != nil {
			failed++
			continue
		}
		err = s.rdb.delViewUV(ctx, id)
		if err != nil {
			failed++
			continue
		}
		success++
	}

	logger.Info(
		"sync article views done",
		zap.Int("total", total),
		zap.Int("success", success),
		zap.Int("failed", failed),
	)
}

type ArticlePageVO struct {
	List     []ArticleSummary `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

func (s *ArticleService) getArticlesByPage(ctx context.Context, page, pageSize int) (*ArticlePageVO, error) {
	var data *ArticlePageVO
	data, err := s.rdb.getArticlesByPage(ctx, page, pageSize)
	if err == nil {
		return data, nil
	}
	if !errors.Is(err, zerrors.CacheMiss) {
		// TODO
	}

	data, err = s.db.listByPage(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	s.rdb.setArticlesByPage(ctx, data)
	return data, nil
}

func (s *ArticleService) getArticlesByPopular(ctx context.Context, limit int) ([]ArticleSummary, error) {
	articles, err := s.rdb.getArticlesByPopular(ctx, limit)
	if err == nil {
		return articles, nil
	}
	if !errors.Is(err, zerrors.CacheMiss) {
		// 这里可能需要输出缓存异常日志
	}

	articles, err = s.db.listPopular(ctx, limit)
	if err != nil {
		return nil, err
	}

	s.rdb.setArticlesByPopular(ctx, limit, articles)
	return articles, nil
}

type CommentVO struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
	Like      string    `json:"like"`
}

type CommentPageVO struct {
	Comments []CommentVO `json:"comments"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `josn:"page_size"`
}

type ArticleDetailVO struct {
	Article  *Article       `json:"article"`
	Comments *CommentPageVO `json:"comments"`
}

func (s *ArticleService) getArticleDetailWithComments(ctx context.Context, id int, userID string, page, pageSize int) (*ArticleDetailVO, error) {
	var articleDetail ArticleDetailVO
	article, err := s.getArticleByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	articleDetail.Article = article

	comments, err := s.getArticleCommentsByPage(ctx, id, page, pageSize)
	if err != nil {
		return nil, err
	}
	articleDetail.Comments = comments

	return &articleDetail, err
}

func (s *ArticleService) getArticleByID(ctx context.Context, id int, userID string) (*Article, error) {
	data, err := s.rdb.getArticleByID(ctx, id)
	if err == nil {
		s.rdb.addViewUV(ctx, id, userID)
		return data, err
	}
	if !errors.Is(err, zerrors.CacheMiss) {
		// 这里可能需要输出缓存异常日志
	}

	data, err = s.db.getArticleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}
	if article == nil {
		return nil, nil
	}

	s.rdb.setArticleByID(ctx, id, data)
	s.rdb.addViewUV(ctx, id, userID)

	return data, nil
}

func (s *ArticleService) getArticleCommentsByPage(ctx context.Context, id int, page, pageSize int) (*CommentPageVO, error) {
	data, err := s.rdb.getArticleComments(ctx, id, page, pageSize)
	if err == nil {
		return data, err
	}
	if !errors.Is(err, zerrors.CacheMiss) {
	}

	var status = Approved
	data, err = s.db.listComment(ctx, CommentQuery{
		status:   &status,
		id:       id,
		page:     page,
		pageSize: pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	s.rdb.setArticleComments(ctx, id, data)
	return data, nil
}

type SendCommentReq struct {
	ArticleID int
	Username  string
	UserIP    string
	Content   string
}

func (s *ArticleService) sendComment(ctx context.Context, req SendCommentReq) (int, error) {
	id, err := s.db.saveComment(ctx, &ArticleComment{
		ArticleID: req.ArticleID,
		Username:  req.Username,
		UserID:    req.UserIP,
		Content:   req.Content,
		Like:      0,
		Status:    Pending,
	})
	if err != nil {
		return -1, fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	return id, err
}

func (s *ArticleService) reviewComment(ctx context.Context, id int, status CommentStatus) error {
	err := s.db.updateCommentStatus(ctx, id, status)
	if err != nil {
		return fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}
	return nil
}

func (s *ArticleService) save(ctx context.Context, article *Article) (int, error) {
	return s.db.save(ctx, article)
}

func (s *ArticleService) delete(ctx context.Context, id int) error {
	return s.db.delete(ctx, id)
}
