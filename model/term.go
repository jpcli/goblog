package model

// 项类型
type TermType uint8

// 项可选值
const (
	TERM_TYPE_CATEGORY TermType = 1
	TERM_TYPE_TAG      TermType = 2
)

// 项表结构
type Term struct {
	Tid         uint32   `db:"tid"`
	Name        string   `db:"name"`
	Slug        string   `db:"slug"`
	Type        TermType `db:"type"`
	Description string   `db:"description"`
	Count       uint32   `db:"count"`
}
