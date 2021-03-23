package repository

import (
	"fmt"

	"github.com/spf13/cast"
)

// InsertPost inserts a new term into 'terms' table.
func (m *mysql) InsertTerm(src *Term) (tid uint32) {
	query := `INSERT INTO terms(name, slug, type, description, count) VALUES(?, ?, ?, ?, ?)`
	_, id := m.Exec(query, src.Name, src.Slug, src.Type, src.Description, src.Count)
	return cast.ToUint32(id)
}

// UpdateTerm updates name and description of term.
func (m *mysql) UpdateTerm(tid uint32, src *Term) {
	query := `UPDATE terms SET name = ?, description = ? WHERE tid = ?`
	m.Exec(query, src.Name, src.Description, tid)
}

func (m *mysql) GetTermByTid(fields string, tid uint32) *Term {
	var ret Term
	query := fmt.Sprintf("SELECT %s FROM terms WHERE tid = ?", fields)
	m.Get(&ret, query, tid)
	return &ret
}

func (m *mysql) GetTermBySlug(fields, slug string) *Term {
	var ret Term
	query := fmt.Sprintf("SELECT %s FROM terms WHERE slug = ?", fields)
	m.Get(&ret, query, slug)
	return &ret
}

func (m *mysql) GetTermList(fields, termType string, desc bool, offset, count uint32) []Term {
	var ret []Term
	var query string
	if desc {
		query = fmt.Sprintf("SELECT %s FROM terms WHERE type = ? ORDER BY tid DESC LIMIT ?, ?", fields)
	} else {
		query = fmt.Sprintf("SELECT %s FROM terms WHERE type = ? LIMIT ?, ?", fields)
	}
	m.Select(&ret, query, termType, offset, count)
	return ret
}

// GetTermsByPid gets one post's all terms by pid and termType.
func (m *mysql) GetTermListByPid(fields, termType string, pid uint32) []Term {
	var ret []Term
	query := fmt.Sprintf(`SELECT %s FROM terms
	WHERE tid in (
		SELECT tid FROM relationships WHERE pid = ?
	) and type = ?`, fields)
	m.Select(&ret, query, pid, termType)
	return ret
}

func (m *mysql) CountTerm(termType string) (count uint32) {
	query := "SELECT COUNT(*) FROM terms WHERE type = ?"
	m.Get(&count, query, termType)
	return
}

// IncrTermsCountByPid increases one post's all terms' count by 1.
func (m *mysql) IncrTermsCountByPid(pid uint32) {
	query := `UPDATE terms SET count = count + 1 
	WHERE tid in (
		SELECT tid FROM relationships WHERE pid = ?
	)`
	m.Exec(query, pid)
}

// DecrTermsCountByPid decreases one post's all terms' count by 1.
func (m *mysql) DecrTermsCountByPid(pid uint32) {
	query := `UPDATE terms SET count = count - 1 
	WHERE tid in (
		SELECT tid FROM relationships WHERE pid = ?
	)`
	m.Exec(query, pid)
}
