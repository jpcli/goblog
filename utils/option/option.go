package option

// TODO: 结合数据库、cache进行返回
func EachPageLimit() uint32 {
	return 10
}

func PageNavLimit() uint32 {
	return 7
}

func WebsiteURL() string {
	return "https://www.jpcli.top"
}
