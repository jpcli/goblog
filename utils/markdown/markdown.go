package markdown

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/russross/blackfriday/v2"
)

func ToHTML(str string) string {
	// blackfriday can't recognize "\r\n" but can recognize "\n"
	text := bytes.Replace([]byte(str), []byte("\r\n"), []byte("\n"), -1)
	// htmlFlags
	r := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CommonHTMLFlags |
			blackfriday.NofollowLinks |
			blackfriday.NoreferrerLinks |
			blackfriday.NoopenerLinks |
			blackfriday.HrefTargetBlank |
			blackfriday.TOC,
	})
	// change h1~h6 to h2~h7
	s := string(blackfriday.Run(text, blackfriday.WithRenderer(r)))
	for i := 6; i > 0; i-- {
		s = strings.Replace(s, fmt.Sprintf("<h%d", i), fmt.Sprintf("<h%d", i+1), -1)
		s = strings.Replace(s, fmt.Sprintf("</h%d>", i), fmt.Sprintf("</h%d>", i+1), -1)
	}
	// add id to the <nav>
	s = strings.Replace(s, "<nav>", `<nav id="post-toc">`, -1)
	// TODO: 删除以下优化，交给前端
	// using bootstrap css to beautify the table
	s = strings.Replace(s, "<table>", `<table class="table table-bordered table-hover table-striped">`, -1)

	return s
}

func ToExcerpt(str string) string {
	s := ToHTML(str)

	reg := regexp.MustCompile(`<pre><code.+?>[\S\s]+?</code></pre>`)
	s = reg.ReplaceAllString(s, "")

	reg = regexp.MustCompile(`<script.+?>[\S\s]+?</script>`)
	s = reg.ReplaceAllString(s, "")

	reg = regexp.MustCompile(`<nav.+?>[\S\s]+?</nav>`)
	s = reg.ReplaceAllString(s, "")

	reg = regexp.MustCompile(`<[^>]+?>`)
	s = reg.ReplaceAllString(s, "")

	reg = regexp.MustCompile(`\s`)
	s = reg.ReplaceAllString(s, ``)

	r := []rune(s)
	if len(r) > 100 {
		return string(r[:100])
	} else {
		return string(r)
	}
}
