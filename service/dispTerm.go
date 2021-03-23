package service

import (
	"goblog/repository"
	"math"
)

const (
	dispTermFields   = "tid, name, type, description, count"
	dispTermliFields = "name, slug, count"
)

// DispTerm is used in post list of term of management system.
type DispTerm struct {
	Tid         uint32 `json:"tid"` // used to get post list by tid
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description" `
	Count       uint32 `json:"count"` // used to generate pageNav
}

// DispTermli means term list item used in frame to show all categories and tags of management system.
type DispTermli struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Count uint32 `json:"count"`
}

func (s *Service) toDispTermList(src []repository.Term) []DispTermli {
	ret := make([]DispTermli, len(src))
	for k, v := range src {
		ret[k] = DispTermli{
			Name:  v.Name,
			Slug:  v.Slug,
			Count: v.Count,
		}
	}
	return ret
}

func (s *Service) GetAllDispCategories() []DispTermli {
	terms := s.repository.GetTermList(dispTermliFields, repository.TermTypeCategory, false, 0, math.MaxUint32)
	return s.toDispTermList(terms)
}

func (s *Service) GetAllDispTags() []DispTermli {
	terms := s.repository.GetTermList(dispTermliFields, repository.TermTypeTag, false, 0, math.MaxUint32)
	return s.toDispTermList(terms)
}
