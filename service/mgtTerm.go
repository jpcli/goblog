package service

import (
	"goblog/repository"
	"goblog/utils/errors"
)

const (
	mgtTermFields   = "name, slug, type, description"
	mgtTermliFields = "tid, name, slug, count"
)

// MgtTerm is used in term edit page of management system.
type MgtTerm struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Type        string `json:"type"`
	Description string `json:"description" `
}

// MgtTermli means term list item used in term list page of management system.
type MgtTermli struct {
	Tid   uint32 `json:"tid"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Count uint32 `json:"count"`
}

// EditTerm creates a new term and returns the new tid when tid=0,
// otherwise updates the term by tid.
func (s *Service) EditTerm(tid uint32, src *MgtTerm) (id uint32, err error) {
	s.repository.Begin()
	if tid == 0 {
		id = s.repository.InsertTerm(&repository.Term{
			Name:        src.Name,
			Slug:        src.Slug,
			Type:        src.Type,
			Description: src.Description,
			Count:       0,
		})
	} else {
		id = tid
		s.repository.UpdateTerm(id, &repository.Term{
			Name:        src.Name,
			Description: src.Description,
		})
	}
	s.repository.Commit()

	err = s.repository.GetError()
	return
}

func (s *Service) GetMgtTermByTid(tid uint32) (*MgtTerm, error) {
	term := s.repository.GetTermByTid(mgtTermFields, tid)
	ret := MgtTerm{
		Name:        term.Name,
		Slug:        term.Slug,
		Type:        term.Type,
		Description: term.Description,
	}

	return &ret, s.repository.GetError()
}

func (s *Service) GetMgtTermList(termType string, page, limit uint32) ([]MgtTermli, error) {
	offset := (page - 1) * limit
	terms := s.repository.GetTermList(mgtTermliFields, termType, true, offset, limit)
	ret := make([]MgtTermli, len(terms))
	for k, v := range terms {
		ret[k] = MgtTermli{
			Tid:   v.Tid,
			Name:  v.Name,
			Slug:  v.Slug,
			Count: v.Count,
		}
	}

	if len(ret) == 0 {
		return nil, errors.WrapErrorWithStack(ErrNoRow)
	}
	return ret, s.repository.GetError()
}

func (s *Service) CountTerm(termType string) (uint32, error) {
	count := s.repository.CountTerm(termType)
	return count, s.repository.GetError()
}
