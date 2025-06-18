package article

import (
	"context"
	"errors"
	"my_web/backend/internal/logger"
	"reflect"
	"testing"
)

type fakeArticleRepo struct {
	getAllArticleIDsFunc     func(ctx context.Context) ([]int, error)
	getArticlesByPageFunc    func(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error)
	getArticleByIDFunc       func(ctx context.Context, id int) (*Article, error)
	getArticlesByPopularFunc func(ctx context.Context, limit int) ([]ArticleWithoutContent, error)
	incrementViewsFunc       func(ctx context.Context, id int, inc int64) error
}

func (f *fakeArticleRepo) getAllArticleIDs(ctx context.Context) ([]int, error) {
	if f.getAllArticleIDsFunc != nil {
		return f.getAllArticleIDsFunc(ctx)
	}
	return nil, nil
}

func (f *fakeArticleRepo) getArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error) {
	if f.getArticlesByPageFunc != nil {
		return f.getArticlesByPageFunc(ctx, page, pageSize)
	}
	return nil, 0, nil
}

func (f *fakeArticleRepo) getArticleByID(ctx context.Context, id int) (*Article, error) {
	if f.getArticleByIDFunc != nil {
		return f.getArticleByIDFunc(ctx, id)
	}
	return nil, nil
}

func (f *fakeArticleRepo) getArticlesByPopular(ctx context.Context, limit int) ([]ArticleWithoutContent, error) {
	if f.getArticlesByPopularFunc != nil {
		return f.getArticlesByPopularFunc(ctx, limit)
	}
	return nil, nil
}

func (f *fakeArticleRepo) incrementViews(ctx context.Context, id int, inc int64) error {
	if f.incrementViewsFunc != nil {
		return f.incrementViewsFunc(ctx, id, inc)
	}
	return nil
}

type fakeArticleCache struct {
	getArticlesByPageFunc    func(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error)
	setArticlesByPageFunc    func(ctx context.Context, page, pageSize int, articles []ArticleWithoutContent, total int) error
	getArticleByIDFunc       func(ctx context.Context, id int) (*Article, error)
	setArticleByIDFunc       func(ctx context.Context, id int, article *Article) error
	getArticlesByPopularFunc func(ctx context.Context, limit int) ([]ArticleWithoutContent, error)
	setArticlesByPopularFunc func(ctx context.Context, limit int, articles []ArticleWithoutContent) error
	addViewUVFunc            func(ctx context.Context, id int, userID string) error
	getViewUVFunc            func(ctx context.Context, id int) (int64, error)
	delViewUVFunc            func(ctx context.Context, id int) error
}

func (f *fakeArticleCache) getArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error) {
	if f.getArticlesByPageFunc != nil {
		return f.getArticlesByPageFunc(ctx, page, pageSize)
	}
	return nil, 0, nil
}

func (f *fakeArticleCache) setArticlesByPage(ctx context.Context, page, pageSize int, articles []ArticleWithoutContent, total int) error {
	if f.setArticlesByPageFunc != nil {
		return f.setArticlesByPageFunc(ctx, page, pageSize, articles, total)
	}
	return nil
}

func (f *fakeArticleCache) getArticleByID(ctx context.Context, id int) (*Article, error) {
	if f.getArticleByIDFunc != nil {
		return f.getArticleByIDFunc(ctx, id)
	}
	return nil, nil
}

func (f *fakeArticleCache) setArticleByID(ctx context.Context, id int, article *Article) error {
	if f.setArticleByIDFunc != nil {
		return f.setArticleByIDFunc(ctx, id, article)
	}
	return nil
}

func (f *fakeArticleCache) getArticlesByPopular(ctx context.Context, limit int) ([]ArticleWithoutContent, error) {
	if f.getArticlesByPopularFunc != nil {
		return f.getArticlesByPopularFunc(ctx, limit)
	}
	return nil, nil
}

func (f *fakeArticleCache) setArticlesByPopular(ctx context.Context, limit int, articles []ArticleWithoutContent) error {
	if f.setArticlesByPopularFunc != nil {
		return f.setArticlesByPopularFunc(ctx, limit, articles)
	}
	return nil
}

func (f *fakeArticleCache) addViewUV(ctx context.Context, id int, userID string) error {
	if f.addViewUVFunc != nil {
		return f.addViewUVFunc(ctx, id, userID)
	}
	return nil
}

func (f *fakeArticleCache) getViewUV(ctx context.Context, id int) (int64, error) {
	if f.getViewUVFunc != nil {
		return f.getViewUVFunc(ctx, id)
	}
	return 0, nil
}

func (f *fakeArticleCache) delViewUV(ctx context.Context, id int) error {
	if f.delViewUVFunc != nil {
		return f.delViewUVFunc(ctx, id)
	}
	return nil
}

func TestArticleService_getArticlesByPage_CacheHit(t *testing.T) {
	logger.InitLogger()
	ctx := context.Background()
	wantArticles := []ArticleWithoutContent{
		{Title: "a"}, {Title: "b"},
	}
	wantTotal := 10
	repo := &fakeArticleRepo{}
	cache := &fakeArticleCache{
		getArticlesByPageFunc: func(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error) {
			if page != 1 || pageSize != 10 {
				t.Fatalf("unexpected page arguments: %d, %d", page, pageSize)
			}
			return wantArticles, wantTotal, nil
		},
	}
	svc := &ArticleService{
		db:  repo,
		rdb: cache,
		cfg: func() *Config { return &Config{} },
	}
	articles, total, err := svc.getArticlesByPage(ctx, 1, 10)
	if err != nil {
		t.Fatalf("getArticlesByPage returned error: %v", err)
	}
	if total != wantTotal {
		t.Fatalf("total = %d, want %d", total, wantTotal)
	}
	if !reflect.DeepEqual(articles, wantArticles) {
		t.Fatalf("articles = %+v, want %+v", articles, wantArticles)
	}
}

func TestArticleService_getArticlesByPage_CacheMiss(t *testing.T) {
	logger.InitLogger()
	ctx := context.Background()
	wantArticles := []ArticleWithoutContent{
		{Title: "c"}, {Title: "d"},
	}
	wantTotal := 20
	cache := &fakeArticleCache{
		getArticlesByPageFunc: func(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error) {
			return nil, 0, ErrCacheMiss
		},
	}
	calledRepo := false
	repo := &fakeArticleRepo{
		getArticlesByPageFunc: func(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error) {
			calledRepo = true
			return wantArticles, wantTotal, nil
		},
	}
	svc := &ArticleService{
		db:  repo,
		rdb: cache,
		cfg: func() *Config { return &Config{} },
	}
	articles, total, err := svc.getArticlesByPage(ctx, 2, 10)
	if err != nil {
		t.Fatalf("getArticlesByPage returned error: %v", err)
	}
	if !calledRepo {
		t.Fatalf("expected repo.getArticlesByPage to be called")
	}
	if total != wantTotal {
		t.Fatalf("total = %d, want %d", total, wantTotal)
	}
	if !reflect.DeepEqual(articles, wantArticles) {
		t.Fatalf("articles = %+v, want %+v", articles, wantArticles)
	}
}

func TestArticleService_getArticlesByPage_CacheError(t *testing.T) {
	logger.InitLogger()
	ctx := context.Background()
	wantArticles := []ArticleWithoutContent{
		{Title: "e"},
	}
	wantTotal := 5
	cache := &fakeArticleCache{
		getArticlesByPageFunc: func(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error) {
			return nil, 0, errors.New("cache error")
		},
	}
	calledRepo := false
	repo := &fakeArticleRepo{
		getArticlesByPageFunc: func(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error) {
			calledRepo = true
			return wantArticles, wantTotal, nil
		},
	}
	svc := &ArticleService{
		db:  repo,
		rdb: cache,
		cfg: func() *Config { return &Config{} },
	}
	articles, total, err := svc.getArticlesByPage(ctx, 1, 1)
	if err != nil {
		t.Fatalf("getArticlesByPage returned error: %v", err)
	}
	if !calledRepo {
		t.Fatalf("expected repo.getArticlesByPage to be called")
	}
	if total != wantTotal {
		t.Fatalf("total = %d, want %d", total, wantTotal)
	}
	if !reflect.DeepEqual(articles, wantArticles) {
		t.Fatalf("articles = %+v, want %+v", articles, wantArticles)
	}
}

func TestArticleService_getArticlesByPopular_CacheHit(t *testing.T) {
	logger.InitLogger()
	ctx := context.Background()
	wantArticles := []ArticleWithoutContent{
		{Title: "hot"},
	}
	cache := &fakeArticleCache{
		getArticlesByPopularFunc: func(ctx context.Context, limit int) ([]ArticleWithoutContent, error) {
			return wantArticles, nil
		},
	}
	repo := &fakeArticleRepo{}
	svc := &ArticleService{
		db:  repo,
		rdb: cache,
		cfg: func() *Config { return &Config{} },
	}
	articles, err := svc.getArticlesByPopular(ctx, 10)
	if err != nil {
		t.Fatalf("getArticlesByPopular returned error: %v", err)
	}
	if !reflect.DeepEqual(articles, wantArticles) {
		t.Fatalf("articles = %+v, want %+v", articles, wantArticles)
	}
}

func TestArticleService_getArticlesByPopular_CacheMiss(t *testing.T) {
	logger.InitLogger()
	ctx := context.Background()
	wantArticles := []ArticleWithoutContent{
		{Title: "dbhot"},
	}
	cache := &fakeArticleCache{
		getArticlesByPopularFunc: func(ctx context.Context, limit int) ([]ArticleWithoutContent, error) {
			return nil, ErrCacheMiss
		},
	}
	calledRepo := false
	repo := &fakeArticleRepo{
		getArticlesByPopularFunc: func(ctx context.Context, limit int) ([]ArticleWithoutContent, error) {
			calledRepo = true
			return wantArticles, nil
		},
	}
	svc := &ArticleService{
		db:  repo,
		rdb: cache,
		cfg: func() *Config { return &Config{} },
	}
	articles, err := svc.getArticlesByPopular(ctx, 3)
	if err != nil {
		t.Fatalf("getArticlesByPopular returned error: %v", err)
	}
	if !calledRepo {
		t.Fatalf("expected repo.getArticlesByPopular to be called")
	}
	if !reflect.DeepEqual(articles, wantArticles) {
		t.Fatalf("articles = %+v, want %+v", articles, wantArticles)
	}
}

func TestArticleService_getArticlesByPopular_CacheError(t *testing.T) {
	logger.InitLogger()
	ctx := context.Background()
	wantArticles := []ArticleWithoutContent{
		{Title: "fallback"},
	}
	cache := &fakeArticleCache{
		getArticlesByPopularFunc: func(ctx context.Context, limit int) ([]ArticleWithoutContent, error) {
			return nil, errors.New("cache error")
		},
	}
	calledRepo := false
	repo := &fakeArticleRepo{
		getArticlesByPopularFunc: func(ctx context.Context, limit int) ([]ArticleWithoutContent, error) {
			calledRepo = true
			return wantArticles, nil
		},
	}
	svc := &ArticleService{
		db:  repo,
		rdb: cache,
		cfg: func() *Config { return &Config{} },
	}
	articles, err := svc.getArticlesByPopular(ctx, 1)
	if err != nil {
		t.Fatalf("getArticlesByPopular returned error: %v", err)
	}
	if !calledRepo {
		t.Fatalf("expected repo.getArticlesByPopular to be called")
	}
	if !reflect.DeepEqual(articles, wantArticles) {
		t.Fatalf("articles = %+v, want %+v", articles, wantArticles)
	}
}

func TestArticleService_getArticleByID_CacheHit(t *testing.T) {
	logger.InitLogger()
	ctx := context.Background()
	wantArticle := &Article{
		Title: "from_cache",
	}
	var addViewCalled bool
	cache := &fakeArticleCache{
		getArticleByIDFunc: func(ctx context.Context, id int) (*Article, error) {
			return wantArticle, nil
		},
		addViewUVFunc: func(ctx context.Context, id int, userID string) error {
			addViewCalled = true
			return nil
		},
	}
	repo := &fakeArticleRepo{}
	svc := &ArticleService{
		db:  repo,
		rdb: cache,
		cfg: func() *Config { return &Config{} },
	}
	article, err := svc.getArticleByID(ctx, 1, "u1")
	if err != nil {
		t.Fatalf("getArticleByID returned error: %v", err)
	}
	if !reflect.DeepEqual(article, wantArticle) {
		t.Fatalf("article = %+v, want %+v", article, wantArticle)
	}
	if !addViewCalled {
		t.Fatalf("expected cache.addViewUV to be called")
	}
}

func TestArticleService_getArticleByID_CacheMiss(t *testing.T) {
	logger.InitLogger()
	ctx := context.Background()
	wantArticle := &Article{
		Title: "from_repo",
	}
	var addViewCalled bool
	var setArticleCalled bool
	cache := &fakeArticleCache{
		getArticleByIDFunc: func(ctx context.Context, id int) (*Article, error) {
			return nil, ErrCacheMiss
		},
		addViewUVFunc: func(ctx context.Context, id int, userID string) error {
			addViewCalled = true
			return nil
		},
		setArticleByIDFunc: func(ctx context.Context, id int, article *Article) error {
			setArticleCalled = true
			return nil
		},
	}
	repoCalled := false
	repo := &fakeArticleRepo{
		getArticleByIDFunc: func(ctx context.Context, id int) (*Article, error) {
			repoCalled = true
			return wantArticle, nil
		},
	}
	svc := &ArticleService{
		db:  repo,
		rdb: cache,
		cfg: func() *Config { return &Config{} },
	}
	article, err := svc.getArticleByID(ctx, 2, "u2")
	if err != nil {
		t.Fatalf("getArticleByID returned error: %v", err)
	}
	if !repoCalled {
		t.Fatalf("expected repo.getArticleByID to be called")
	}
	if !addViewCalled {
		t.Fatalf("expected cache.addViewUV to be called")
	}
	if !setArticleCalled {
		t.Fatalf("expected cache.setArticleByID to be called")
	}
	if !reflect.DeepEqual(article, wantArticle) {
		t.Fatalf("article = %+v, want %+v", article, wantArticle)
	}
}

func TestArticleService_getArticleByID_CacheError(t *testing.T) {
	logger.InitLogger()
	ctx := context.Background()
	wantArticle := &Article{
		Title: "from_repo_on_error",
	}
	cache := &fakeArticleCache{
		getArticleByIDFunc: func(ctx context.Context, id int) (*Article, error) {
			return nil, errors.New("cache error")
		},
	}
	repoCalled := false
	var addViewCalled bool
	var setArticleCalled bool
	cache.addViewUVFunc = func(ctx context.Context, id int, userID string) error {
		addViewCalled = true
		return nil
	}
	cache.setArticleByIDFunc = func(ctx context.Context, id int, article *Article) error {
		setArticleCalled = true
		return nil
	}
	repo := &fakeArticleRepo{
		getArticleByIDFunc: func(ctx context.Context, id int) (*Article, error) {
			repoCalled = true
			return wantArticle, nil
		},
	}
	svc := &ArticleService{
		db:  repo,
		rdb: cache,
		cfg: func() *Config { return &Config{} },
	}
	article, err := svc.getArticleByID(ctx, 3, "u3")
	if err != nil {
		t.Fatalf("getArticleByID returned error: %v", err)
	}
	if !repoCalled {
		t.Fatalf("expected repo.getArticleByID to be called")
	}
	if !addViewCalled {
		t.Fatalf("expected cache.addViewUV to be called")
	}
	if !setArticleCalled {
		t.Fatalf("expected cache.setArticleByID to be called")
	}
	if !reflect.DeepEqual(article, wantArticle) {
		t.Fatalf("article = %+v, want %+v", article, wantArticle)
	}
}
