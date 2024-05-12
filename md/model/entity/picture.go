package entity

type Picture struct {
	Id         string `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Path       string `json:"path" db:"path"`
	Hash       string `json:"hash" db:"hash"`
	Size       int64  `json:"size" db:"size"`
	CreateTime int64  `json:"createTime" db:"create_time"`
	UserId     string `json:"userId" db:"user_id"`
}

type PicturePageResult struct {
	Picture
	PicturePrefix   string `json:"picturePrefix"`
	ThumbnailPrefix string `json:"thumbnailPrefix"`
}
