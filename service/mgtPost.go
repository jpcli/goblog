package service

import (
	"goblog/repository"
	"goblog/utils/errors"
	"goblog/utils/markdown"
	"math"
	"time"

	"github.com/spf13/cast"
)

const (
	mgtPostFields         = "title, modified, excerpt, keywords, text, commentAllow"
	mgtPostliFields       = "pid, title, modified, status"
	mgtCatgTagsItemFields = "tid, name"
)

// MgtPost is used in post edit page of management system.
type MgtPost struct {
	Title        string            `json:"title"`
	Modified     uint32            `json:"modified"`
	Excerpt      string            `json:"excerpt"`
	Keywords     string            `json:"keywords"`
	Text         string            `json:"text"`
	CommentAllow uint8             `json:"commentAllow"`
	Category     MgtCatgTagsItem   `json:"category"`
	Tags         []MgtCatgTagsItem `json:"tags"`
}

// MgtPostli means post list item used in post list page of management system.
type MgtPostli struct {
	Pid      uint32            `json:"pid"`
	Title    string            `json:"title"`
	Modified uint32            `json:"modified"`
	Status   string            `json:"status"`
	Category MgtCatgTagsItem   `json:"category"`
	Tags     []MgtCatgTagsItem `json:"tags"`
}

type MgtCatgTagsItem struct {
	Tid  uint32 `json:"tid"`
	Name string `json:"name"`
}

func (s *Service) toMgtCatgTagsList(src []repository.Term) []MgtCatgTagsItem {
	ret := make([]MgtCatgTagsItem, len(src))
	for k, v := range src {
		ret[k] = MgtCatgTagsItem{
			Tid:  v.Tid,
			Name: v.Name,
		}
	}
	return ret
}

func (s *Service) GetMgtPostByPid(pid uint32) (*MgtPost, error) {
	if pid == 0 {
		return nil, nil
	}

	post := s.repository.GetPostByPid(mgtPostFields, pid)
	catg := s.repository.GetTermListByPid(mgtCatgTagsItemFields, repository.TermTypeCategory, pid)
	tags := s.repository.GetTermListByPid(mgtCatgTagsItemFields, repository.TermTypeTag, pid)
	if err := s.repository.GetError(); err != nil {
		return nil, err
	}

	ret := MgtPost{
		Title:        post.Title,
		Modified:     post.Modified,
		Excerpt:      post.Excerpt,
		Keywords:     post.Keywords,
		Text:         post.Text,
		CommentAllow: post.CommentAllow,
		Category:     s.toMgtCatgTagsList(catg)[0],
		Tags:         s.toMgtCatgTagsList(tags),
	}
	return &ret, nil
}

// EditPost creates a new post and returns the new pid when pid==0,
// otherwise updates the post by pid.
func (s *Service) EditPost(pid uint32, src *MgtPost) (id uint32, err error) {
	// TODO:校验放在handler层
	if src.Category.Tid == 0 {
		err = errors.NewErrorWithStack("miss category of post")
		return
	}
	//TODO:验证分类、标签中id是否为对应type

	now := cast.ToUint32(time.Now().Unix())
	s.repository.Begin()

	if pid == 0 {
		id = s.repository.InsertPost(&repository.Post{
			Title:        src.Title,
			Created:      now,
			Modified:     now,
			Excerpt:      markdown.ToExcerpt(src.Text),
			Keywords:     src.Keywords,
			Text:         src.Text,
			Status:       repository.PostStatusPublish,
			CommentCount: 0,
			CommentAllow: src.CommentAllow,
		})
		s.repository.InsertPostmeta(&repository.Postmeta{
			Pid:       id,
			MetaKey:   "post_view_count",
			MetaValue: "0",
		})
	} else {
		id = pid
		src.Modified = now
		s.repository.UpdatePost(id, &repository.Post{
			Title:        src.Title,
			Modified:     now,
			Excerpt:      markdown.ToExcerpt(src.Text),
			Keywords:     src.Keywords,
			Text:         src.Text,
			CommentAllow: src.CommentAllow,
		})
		s.repository.DecrTermsCountByPid(id)
		s.repository.DeleteRelationshipsByPid(id)
	}

	// insert new relationships
	args := make([]*repository.Relationship, 1+len(src.Tags))
	args[0] = &repository.Relationship{Pid: id, Tid: src.Category.Tid}
	for k, v := range src.Tags {
		args[k+1] = &repository.Relationship{Pid: id, Tid: v.Tid}
	}
	s.repository.InsertRelationships(args...)
	s.repository.IncrTermsCountByPid(id)

	s.repository.Commit()
	err = s.repository.GetError()
	return
}

func (s *Service) ModifyPostStatus(pid uint32, status string) error {
	s.repository.Begin()
	s.repository.UpdatePostStatus(pid, status)
	s.repository.Commit()
	return s.repository.GetError()
}

// GetMgtPostList will make a slice of len equal to limit first, so don't use MaxUint32 as limit.
func (s *Service) GetMgtPostList(page, limit uint32) ([]MgtPostli, error) {
	result := make([]repository.Post, 0, limit)

	stickyCount := s.repository.CountPost(repository.PostStatusSticky)
	n, m := (page-1)*limit, page*limit
	if page <= uint32(math.Floor(float64(stickyCount)/float64(limit))) {
		// all sticky posts
		result = append(result, s.repository.GetPostListByStatus(mgtPostliFields, repository.PostStatusSticky, n, limit)...)
	} else if page > uint32(math.Ceil(float64(stickyCount)/float64(limit))) {
		// all publish posts
		result = append(result, s.repository.GetPostListByStatus(mgtPostliFields, repository.PostStatusPublish, n-stickyCount, limit)...)
	} else {
		// sticky posts and publish posts mixed
		result = append(result, s.repository.GetPostListByStatus(mgtPostliFields, repository.PostStatusSticky, n, stickyCount-n)...)
		result = append(result, s.repository.GetPostListByStatus(mgtPostliFields, repository.PostStatusPublish, 0, m-stickyCount)...)
	}

	ret := make([]MgtPostli, len(result))
	for k, v := range result {
		catg := s.repository.GetTermListByPid(mgtCatgTagsItemFields, repository.TermTypeCategory, v.Pid)
		tags := s.repository.GetTermListByPid(mgtCatgTagsItemFields, repository.TermTypeTag, v.Pid)
		ret[k] = MgtPostli{
			Pid:      v.Pid,
			Title:    v.Title,
			Modified: v.Modified,
			Status:   v.Status,
			Category: s.toMgtCatgTagsList(catg)[0],
			Tags:     s.toMgtCatgTagsList(tags),
		}
	}

	if len(ret) == 0 {
		return nil, errors.WrapErrorWithStack(ErrNoRow)
	}
	return ret, s.repository.GetError()
}

func (s *Service) CountPost() (uint32, error) {
	stickyCount := s.repository.CountPost(repository.PostStatusSticky)
	publishCount := s.repository.CountPost(repository.PostStatusPublish)
	return stickyCount + publishCount, s.repository.GetError()
}
