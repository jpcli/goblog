package request

type Post struct {
	ID           uint32   `json:"id"`
	Title        string   `json:"title"`
	CateID       uint32   `json:"category_id"`
	TagsID       []uint32 `json:"tags_id"`
	Keywords     string   `json:"keywords"`
	CommentAllow uint32   `json:"comment_allow"`
	Text         string   `json:"text"`
}
