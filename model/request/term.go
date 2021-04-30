package request

type Term struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Type        uint8  `json:"type"`
	Description string `json:"description"`
}

type TermList struct {
	TermType uint8 `json:"term_type" form:"term_type"`
	Pager
}
