package entity

type User struct {
	Id         string `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Password   string `json:"password" db:"password"`
	CreateTime int64  `json:"createTime" db:"create_time"`
}

type UserCondition struct {
	Id          string `json:"id"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}
