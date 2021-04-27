package model

// 选项表结构
type Option struct {
	Name  string `db:"name"`
	Value string `db:"value"`
}
