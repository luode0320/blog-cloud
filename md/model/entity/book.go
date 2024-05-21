package entity

type Book struct {
	Id         string `json:"id" db:"id"`
	ParentId   string `json:"parentId" db:"parent_id"`
	Name       string `json:"name" db:"name"`
	CreateTime int64  `json:"createTime" db:"create_time"`
	UserId     string `json:"userId" db:"user_id"`
}
