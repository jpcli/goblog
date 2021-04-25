package model

// 附件表结构
type Attachment struct {
	Aid          uint32 `db:"aid"`
	Created      uint32 `db:"created"`
	RelativePath string `db:"relativePath"`
}
