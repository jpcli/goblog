package service

import (
	"encoding/xml"
	"fmt"
	"goblog/repository"
	"goblog/utils/option"
	"math"
	"time"
)

type urlset struct {
	XMLName     struct{} `xml:"urlset"`
	URL         []*url   `xml:",innerxml"`
	Xmlns       string   `xml:"xmlns,attr"`
	XmlnsMobile string   `xml:"xmlns:mobile,attr"`
}

type url struct {
	XMLName    struct{}  `xml:"url"`
	Loc        string    `xml:"loc,omitempty"`
	Lastmod    string    `xml:"lastmod,omitempty"`
	Changefreq string    `xml:"changefreq,omitempty"`
	MobileTag  mobileTag `xml:",innerxml"`
}
type mobileTag struct {
	XMLName struct{} `xml:"mobile:mobile"`
	Type    string   `xml:"type,attr"`
}

func newURLset() *urlset {
	return &urlset{
		Xmlns:       "http://www.sitemaps.org/schemas/sitemap/0.9",
		XmlnsMobile: "http://www.baidu.com/schemas/sitemap-mobile/1/",
		URL:         make([]*url, 0, 100),
	}
}

func (s *urlset) addURL(u *url) {
	s.URL = append(s.URL, u)
}

func newURL(loc string) *url {
	return &url{
		Loc:        loc,
		Changefreq: "daily",
		MobileTag: mobileTag{
			Type: `pc,mobile`,
		},
	}
}

func (u *url) setLastmod(lastmod string) *url {
	u.Lastmod = lastmod
	return u
}

// func (u *url) setChangefreq(changefreq string) *url {
// 	u.Changefreq = changefreq
// 	return u
// }

func (s *Service) GetSitemap() (string, *DispError) {
	URLset := newURLset()
	website := option.WebsiteURL()
	pages := []string{"archives", "declaration"}

	posts := s.repository.GetPostList("pid, modified", 0, math.MaxUint32)
	for _, v := range posts {
		u := newURL(fmt.Sprintf("%s/post/%d.html", website, v.Pid)).
			setLastmod(time.Unix(int64(v.Modified), 0).Format("2006-01-02"))
		URLset.addURL(u)
	}

	categories := s.repository.GetTermList("slug", repository.TermTypeCategory, false, 0, math.MaxUint32)
	for _, v := range categories {
		u := newURL(fmt.Sprintf("%s/category/%s/", website, v.Slug))
		URLset.addURL(u)
	}

	tags := s.repository.GetTermList("slug", repository.TermTypeTag, false, 0, math.MaxUint32)
	for _, v := range tags {
		u := newURL(fmt.Sprintf("%s/tag/%s/", website, v.Slug))
		URLset.addURL(u)
	}

	if err := s.repository.GetError(); err != nil {
		return "", NewDispError(ErrorServer, err)
	}

	for _, v := range pages {
		u := newURL(fmt.Sprintf("%s/page/%s.html", website, v))
		URLset.addURL(u)
	}

	d, err := xml.Marshal(URLset)
	if err != nil {
		return "", NewDispError(ErrorServer, err)
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>%s`, string(d)), nil
}
