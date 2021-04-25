package model

// 选项表结构
type Option struct {
	OptionKey   string `db:"optionKey"`
	OptionValue string `db:"optionValue"`
}
