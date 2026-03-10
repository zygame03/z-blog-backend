package article

import (
	"fmt"
)

func articleTotalKey() string {
	return "Article:Total"
}

func articleByIDKey(id int) string {
	return fmt.Sprintf("Article:ByID:%d", id)
}

func articleByPageKey(page, pageSize int) string {
	return fmt.Sprintf("Article:ByPage:%d:%d", page, pageSize)
}

func articleByPopularKey(limit int) string {
	return fmt.Sprintf("Article:ByPopular:%d", limit)
}

func articleActiveViewIDsKey() string {
	return "Article:View:ActiveIDs"
}

func articleViewKey(id int) string {
	if id == -1 {
		return "Article:view:UV:*"
	}
	return fmt.Sprintf("Article:View:UV:%d", id)
}
