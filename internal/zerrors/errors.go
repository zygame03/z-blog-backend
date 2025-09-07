package zerrors

import "errors"

var (
	ErrInternal       = errors.New("未知错误")
	ErrRequest        = errors.New("异常请求")
	ErrCacheMiss      = errors.New("缓存未命中")
	ErrDBOperation    = errors.New("数据库操作异常")
	ErrCacheOperation = errors.New("缓存操作异常")

	ErrArticleNotFound = errors.New("未找到文章")

	ErrUserNotFound     = errors.New("未知用户")
	ErrUserAlreadyExist = errors.New("用户已存在")
	ErrInvalidPassword  = errors.New("密码错误")

	ErrUnmarshal = errors.New("反序列化失败")
	ErrMarshal   = errors.New("序列化失败")
	ErrParse     = errors.New("数据解析失败")
)
