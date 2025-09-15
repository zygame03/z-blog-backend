package site

import "fmt"

func getIntroKey() string {
	return "blog:site:intro"
}

func announcementKey() string {
	return "blog:site:announcement:all"
}

func getAnnouncementKey(id int) string {
	return fmt.Sprintf("blog:site:announcement:%d", id)
}

func danmakuKey() string {
	return "blog:site:danmaku"
}
