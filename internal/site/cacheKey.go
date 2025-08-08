package site

import "fmt"

func getIntroKey() string {
	return "websiteData:intro"
}

func getAllAnnouncementKey() string {
	return "site:announcement:*"
}

func getAnnouncementKey(id int) string {
	return fmt.Sprintf("site:announcement:%d", id)
}
