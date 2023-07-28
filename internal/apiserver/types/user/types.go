package user

import (
	"gitlab.intsig.net/textin-gateway/internal/apiserver/types/base"
)

type EditInfo struct {
	ID         int32 `json:"id"`
	CreateTime int64 `json:"create_time"`
	UpdateTime int64 `json:"update_time"`
}

type UserDef struct {
	Alias    string  `json:"alias"`
	Account  string  `json:"account"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
	Channels []int32 `json:"channels"`
	Desc     string  `json:"desc"`
	Status   *bool   `json:"status"`
}

type User struct {
	EditInfo
	Status   *bool           `json:"status"`
	Alias    string          `json:"alias"`
	Account  string          `json:"account"`
	Role     string          `json:"role"`
	Channels []base.BaseInfo `json:"channels"`
	Desc     string          `json:"desc"`
}

type ListReq struct {
	ID       int    `form:"id"`
	Account  string `form:"account"`
	Alias    string `form:"alias"`
	Role     string `form:"role"`
	Page     int    `form:"page"`      //页数
	PageSize int    `form:"page_size"` //页大小
}

type UserList struct {
	Items []User `json:"items"`
	Count uint64 `json:"count"`
}

type UpdateUser struct {
	ID int32 `uri:"id"`
	UserDef
}

type IdPathParam struct {
	ID int32 `uri:"id"`
}

type ResetPasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type GroupList struct {
	Items []base.BaseInfo `json:"items"`
	Count uint64          `json:"count"`
}

type ListGroupReq struct {
	ID       int    `form:"id"`
	Name     string `form:"name"`
	Page     int    `form:"page"`      //页数
	PageSize int    `form:"page_size"` //页大小
}
