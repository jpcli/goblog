package request

type Pager struct {
	Pi uint32 `json:"page" form:"page"`
	Ps uint32 `json:"limit" form:"limit"`
}
