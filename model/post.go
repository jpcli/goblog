package model

// 文章状态
const (
	POST_STATUS_PUBLISH uint8 = 1
	POST_STATUS_STICKY  uint8 = 2
	POST_STATUS_DELETED uint8 = 3
)

// 文章是否允许评论
const (
	POST_COMMENT_ALLOW    uint8 = 1
	POST_COMMENT_DISALLOW uint8 = 2
)

// 文章表结构
type Post struct {
	Pid          uint32 `db:"pid"`
	Title        string `db:"title"`
	Created      uint32 `db:"created"`
	Modified     uint32 `db:"modified"`
	Excerpt      string `db:"excerpt"`
	Keywords     string `db:"keywords"`
	Text         string `db:"text"`
	Status       uint8  `db:"status"`
	CommentCount uint32 `db:"commentCount"`
	CommentAllow uint8  `db:"commentAllow"`
}

// 文章元数据
type Postmeta struct {
	Mid       uint32 `db:"mid"`
	Pid       uint32 `db:"pid"`
	MetaKey   string `db:"metaKey"`
	MetaValue string `db:"metaValue"`
}
