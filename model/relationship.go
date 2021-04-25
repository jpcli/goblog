package model

// 关系表结构
type Relationship struct {
	Pid uint32 `db:"pid"`
	Tid uint32 `db:"tid"`
}
