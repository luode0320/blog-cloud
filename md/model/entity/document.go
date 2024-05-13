package entity

type Document struct {
	Id         string       `json:"id" db:"id"`
	Name       string       `json:"name" db:"name"`
	Content    string       `json:"content" db:"content"`
	Type       DocumentType `json:"type" db:"type"`
	Published  bool         `json:"published" db:"published"`
	CreateTime int64        `json:"createTime" db:"create_time"`
	UpdateTime int64        `json:"updateTime" db:"update_time"`
	BookId     string       `json:"bookId" db:"book_id"`
	UserId     string       `json:"userId" db:"user_id"`
}

type DocumentPageResult struct {
	Id         string       `json:"id" db:"id"`
	Name       string       `json:"name" db:"name"`
	Type       DocumentType `json:"type" db:"type"`
	CreateTime int64        `json:"createTime" db:"create_time"`
	UpdateTime int64        `json:"updateTime" db:"update_time"`
	Username   string       `json:"username" db:"username"`
	BookName   string       `json:"bookName" db:"book_name"`
}

type DocumentPageCondition struct {
	Username string       `json:"username"`
	Name     string       `json:"name"`
	Type     DocumentType `json:"type"`
	BookName string       `json:"bookName"`
}

type DocumentType string

const (
	MdExt      string       = ".md"     // MD文件后缀
	DocMd      DocumentType = "md"      // 文档类型：markdown
	DocOpenApi DocumentType = "openApi" // 文档类型：OpenApi
)
