package user

import "my_web/backend/internal/global"

type Profile struct {
	global.BaseModel
	Name   string
	Desc   string
	Locate string
	Age    int
}
