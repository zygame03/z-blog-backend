package article

import (
	"fmt"
)

func commentsByIdKey(id int, page, pageSize int) string {
	return fmt.Sprintf("blog:article:comment:%d:page:%d:pageSize:%d", id, page, pageSize)
}

func articleTotalKey() string {
	return "blog:article:total"
}

func articleByIDKey(id int) string {
	return fmt.Sprintf("blog:article:detail:%d", id)
}

func articleByPageKey(page, pageSize int) string {
	return fmt.Sprintf("blog:article:page:%d:%d", page, pageSize)
}

func articleByPopularKey(limit int) string {
	return fmt.Sprintf("blog:article:popular:%d", limit)
}

func articleActiveViewIDsKey() string {
	return "blog:article:view:active"
}

func articleViewKey(id int) string {
	if id == -1 {
		return "blog:article:view:UV:*"
	}
	return fmt.Sprintf("blog:article:view:UV:%d", id)
}
