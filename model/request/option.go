package request

type BaseOption struct {
	PageSize    uint32 `json:"page_size"`
	PageNavSize uint32 `json:"page_nav_size"`
	SiteName    string `json:"site_name"`
	SiteURL     string `json:"site_url"`
}
