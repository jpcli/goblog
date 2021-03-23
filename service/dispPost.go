package service

import (
	"fmt"
	"goblog/repository"
	"goblog/utils/errors"
	"goblog/utils/option"
	"math"
	"strings"
	"time"

	"github.com/spf13/cast"
)

const (
	dispPostFields         = "title, created, modified, excerpt, keywords, text, commentCount, commentAllow"
	dispPostliFields       = "pid, title, created, excerpt, status"
	dispCatgTagsItemFields = "slug, name"
)

// DispPost is used in post show page of display system.
type DispPost struct {
	Title        string `db:"title"`
	Created      uint32 `db:"created"`
	Modified     uint32 `db:"modified"`
	Excerpt      string `db:"excerpt"`
	Keywords     string `db:"keywords"`
	Text         string `db:"text"`
	CommentCount uint32 `db:"commentCount"`
	CommentAllow uint8  `db:"commentAllow"`
	ViewCount    uint32
	Category     DispCatgTagsItem
	Tags         []DispCatgTagsItem
}

// DispPostli means post list item used in post list page of display system.
type DispPostli struct {
	Pid      uint32 `db:"pid"`
	Title    string `db:"title"`
	Created  uint32 `db:"created"`
	Excerpt  string `db:"excerpt"`
	Status   string `db:"status"`
	Category DispCatgTagsItem
}

type DispCatgTagsItem struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func (s *Service) toDispCatgTagsList(src []repository.Term) []DispCatgTagsItem {
	ret := make([]DispCatgTagsItem, len(src))
	for k, v := range src {
		ret[k] = DispCatgTagsItem{
			Slug: v.Slug,
			Name: v.Name,
		}
	}
	return ret
}

func (s *Service) toDispPostList(posts []repository.Post) ([]DispPostli, *DispError) {
	ret := make([]DispPostli, len(posts))
	for k, v := range posts {
		catg := s.repository.GetTermListByPid(dispCatgTagsItemFields, repository.TermTypeCategory, v.Pid)
		ret[k] = DispPostli{
			Pid:      v.Pid,
			Title:    v.Title,
			Created:  v.Created,
			Excerpt:  v.Excerpt,
			Status:   v.Status,
			Category: s.toDispCatgTagsList(catg)[0],
		}
	}

	if err := s.repository.GetError(); err != nil {
		return nil, NewDispError(ErrorServer, err)
	} else if len(ret) == 0 {
		return nil, NewDispError(ErrorRequest, errors.NewErrorWithStack("no result"))
	}

	return ret, nil
}

func (s *Service) GetDispPostByPid(pid uint32) (*DispPost, *DispError) {
	post := s.repository.GetPostByPid(dispPostFields, pid)
	catg := s.repository.GetTermListByPid(dispCatgTagsItemFields, repository.TermTypeCategory, pid)
	tags := s.repository.GetTermListByPid(dispCatgTagsItemFields, repository.TermTypeTag, pid)
	if err := s.repository.GetError(); err != nil {
		if errors.Cause(err) == repository.ErrNoRow {
			return nil, NewDispError(ErrorRequest, err)
		} else {
			return nil, NewDispError(ErrorServer, err)
		}
	}

	viewCount := cast.ToUint32(s.repository.GetPostmetaValue(pid, "post_view_count")) + 1

	if err := s.repository.GetError(); err != nil {
		return nil, NewDispError(ErrorServer, err)
	}

	return &DispPost{
		Title:        post.Title,
		Created:      post.Created,
		Modified:     post.Modified,
		Excerpt:      post.Excerpt,
		Keywords:     post.Keywords,
		Text:         post.Text,
		CommentCount: post.CommentCount,
		CommentAllow: post.CommentAllow,
		ViewCount:    viewCount,
		Category:     s.toDispCatgTagsList(catg)[0],
		Tags:         s.toDispCatgTagsList(tags),
	}, nil
}

func (s *Service) IncrPostViewCountByPid(pid uint32) error {
	s.repository.Begin()
	s.repository.IncrPostViewCountByPid(pid)
	s.repository.Commit()
	return s.repository.GetError()
}

// GetDispPostList will make a slice of len equal to limit first, so don't use MaxUint32 as limit.
func (s *Service) GetDispPostList(page, limit uint32) ([]DispPostli, *PageNav, *DispError) {
	// TODO:重构，取result的代码与GetMgtPostList中代码重复
	result := make([]repository.Post, 0, limit)

	stickyCount := s.repository.CountPost(repository.PostStatusSticky)
	n, m := (page-1)*limit, page*limit
	if page <= uint32(math.Floor(float64(stickyCount)/float64(limit))) {
		// all sticky posts
		result = append(result, s.repository.GetPostListByStatus(dispPostliFields, repository.PostStatusSticky, n, limit)...)
	} else if page > uint32(math.Ceil(float64(stickyCount)/float64(limit))) {
		// all publish posts
		result = append(result, s.repository.GetPostListByStatus(dispPostliFields, repository.PostStatusPublish, n-stickyCount, limit)...)
	} else {
		// sticky posts and publish posts mixed
		result = append(result, s.repository.GetPostListByStatus(dispPostliFields, repository.PostStatusSticky, n, stickyCount-n)...)
		result = append(result, s.repository.GetPostListByStatus(dispPostliFields, repository.PostStatusPublish, 0, m-stickyCount)...)
	}

	posts, e := s.toDispPostList(result)
	if e != nil {
		return nil, nil, e
	}

	count, err := s.CountPost()
	if err != nil {
		return nil, nil, NewDispError(ErrorServer, err)
	}

	pageNav, e := getPageNav(count, page, limit, option.PageNavLimit())

	return posts, pageNav, e
}

func (s *Service) GetDispPostListBySlug(metaType, slug string, page, limit uint32) ([]DispPostli, *PageNav, *repository.Term, *DispError) {
	t := s.repository.GetTermBySlug(dispTermFields, slug)
	if err := s.repository.GetError(); err != nil && errors.Cause(err) != repository.ErrNoRow {
		return nil, nil, nil, NewDispError(ErrorServer, err)
	} else if t.Type != metaType || errors.Cause(err) == repository.ErrNoRow {
		return nil, nil, nil, NewDispError(ErrorRequest, fmt.Errorf("no such %s: %s", metaType, slug))
	}

	offset := (page - 1) * limit
	p := s.repository.GetPostListByTid(dispPostliFields, t.Tid, offset, limit)
	posts, e := s.toDispPostList(p)
	if e != nil {
		return nil, nil, nil, e
	}

	pageNav, e := getPageNav(t.Count, page, limit, option.PageNavLimit())

	return posts, pageNav, t, e
}

type PageNav struct {
	TotalPage uint32
	Current   uint32
	Previous  uint32
	Next      uint32
	PageList  []uint32
	URLFormat string
}

func getPageNav(postCount, page, limit, displayNum uint32) (*PageNav, *DispError) {
	if displayNum%2 == 0 {
		return nil, NewDispError(ErrorServer, errors.NewErrorWithStack("display should be an odd number"))
	}

	totalPage := uint32(math.Ceil(float64(postCount) / float64(limit)))

	// generate the pageList
	var list []uint32
	mid := (displayNum - 1) / 2
	if totalPage < displayNum {
		// 1 ~ total
		list = make([]uint32, totalPage)
		for i := uint32(0); i < totalPage; i++ {
			list[i] = i + 1
		}
	} else {
		list = make([]uint32, displayNum)
		var begin uint32
		if page <= mid {
			begin = 1 // 1 ~ displayNum
		} else if page > totalPage-mid {
			begin = totalPage - displayNum + 1 // totalPage-displayNum+1 ~ totalPage
		} else {
			begin = page - mid // page-mid ~ page+mid
		}
		for i := uint32(0); i < displayNum; i++ {
			list[i] = begin + 1
		}
	}

	// generate the pageNav info
	return &PageNav{
		TotalPage: totalPage,
		Current:   page,
		Previous:  page - 1,
		Next:      page + 1,
		PageList:  list,
	}, nil

}

func (s *Service) GetArchives() ([]map[string]interface{}, *DispError) {
	posts := s.repository.GetPostList("pid, title, created", 0, math.MaxUint32)
	if err := s.repository.GetError(); err != nil {
		return nil, NewDispError(ErrorServer, err)
	}

	type ArchiveItem struct {
		Pid   uint32
		Title string
		Day   string
	}

	ret := make([]map[string]interface{}, 0, 10)
	var previous string
	nowIndex := -1
	for _, v := range posts {
		date := strings.Split(time.Unix(int64(v.Created), 0).Format("2006年01月 02日"), " ")
		// new group
		if date[0] != previous {
			previous = date[0]
			nowIndex++
			// append a new group to the slice
			ret = append(ret, make(map[string]interface{}))
			// init a new group
			ret[nowIndex]["groupName"] = date[0]
			ret[nowIndex]["groupItems"] = make([]ArchiveItem, 0, 10)
		}
		// add the item to the group
		ret[nowIndex]["groupItems"] = append(ret[nowIndex]["groupItems"].([]ArchiveItem), ArchiveItem{
			Pid:   v.Pid,
			Title: v.Title,
			Day:   date[1],
		})
	}
	return ret, nil
}

func (s *Service) GetLastModified() uint32 {
	return s.repository.GetLastModified()
}
