package repository

// the value range of status field in 'posts' table.
const (
	PostStatusPublish = "publish"
	PostStatusSticky  = "sticky"
	PostStatusTrash   = "trash"
)

// Post is the table field of posts.
type Post struct {
	Pid          uint32 `db:"pid"`
	Title        string `db:"title"`
	Created      uint32 `db:"created"`
	Modified     uint32 `db:"modified"`
	Excerpt      string `db:"excerpt"`
	Keywords     string `db:"keywords"`
	Text         string `db:"text"`
	Status       string `db:"status"`
	CommentCount uint32 `db:"commentCount"`
	CommentAllow uint8  `db:"commentAllow"`
}

// Postmeta is the table field of postmeta.
type Postmeta struct {
	Mid       uint32 `db:"mid"`
	Pid       uint32 `db:"pid"`
	MetaKey   string `db:"metaKey"`
	MetaValue string `db:"metaValue"`
}

// the value range of type field in 'terms' table.
const (
	TermTypeCategory = "category"
	TermTypeTag      = "tag"
)

// Term is the table field of terms.
type Term struct {
	Tid         uint32 `db:"tid"`
	Name        string `db:"name"`
	Slug        string `db:"slug"`
	Type        string `db:"type"`
	Description string `db:"description"`
	Count       uint32 `db:"count"`
}

// Relationship is the table field of relationships.
type Relationship struct {
	Pid uint32 `db:"pid"`
	Tid uint32 `db:"tid"`
}

// Attachment is the table field of attachments.
type Attachment struct {
	Aid          uint32 `db:"aid"`
	Created      uint32 `db:"created"`
	RelativePath string `db:"relativePath"`
}
