package cache

import (
	"errors"
	"fmt"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

func ArticleTotal() string {
	return "Article:Total"
}

func ArticleByIDKey(id int) string {
	return fmt.Sprintf("Article:ByID:%d", id)
}

func ArticleByPageKey(page, pageSize int) string {
	return fmt.Sprintf("Article:ByPage:%d:%d", page, pageSize)
}

func ArticleByPopularKey(limit int) string {
	return fmt.Sprintf("Article:ByPopular:%d", limit)
}
