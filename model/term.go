package model

// 项类型
const (
	TERM_TYPE_CATEGORY uint8 = 1
	TERM_TYPE_TAG      uint8 = 2
)

// 项表结构
type Term struct {
	Tid         uint32 `db:"tid"`
	Name        string `db:"name"`
	Slug        string `db:"slug"`
	Type        uint8  `db:"type"`
	Description string `db:"description"`
	Count       uint32 `db:"count"`
}
