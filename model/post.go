package model

// 文章状态类型
type PostStatusType uint8

// 文章状态可选值
const (
	POST_STATUS_PUBLISH PostStatusType = 1
	POST_STATUS_STICKY  PostStatusType = 2
	POST_STATUS_DELETED PostStatusType = 3
)

// 文章是否允许评论类型
type PostCommentAllowType uint8

// 文章是否允许评论可选值
const (
	POST_COMMENT_ALLOW    PostCommentAllowType = 1
	POST_COMMENT_DISALLOW PostCommentAllowType = 2
)

// 文章表结构
type Post struct {
	Pid          uint32               `db:"pid"`
	Title        string               `db:"title"`
	Created      uint32               `db:"created"`
	Modified     uint32               `db:"modified"`
	Excerpt      string               `db:"excerpt"`
	Keywords     string               `db:"keywords"`
	Text         string               `db:"text"`
	Status       PostStatusType       `db:"status"`
	CommentCount uint32               `db:"commentCount"`
	CommentAllow PostCommentAllowType `db:"commentAllow"`
}

// 文章元数据
type Postmeta struct {
	Mid       uint32 `db:"mid"`
	Pid       uint32 `db:"pid"`
	MetaKey   string `db:"metaKey"`
	MetaValue string `db:"metaValue"`
}
