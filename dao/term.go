package dao

import (
	"database/sql"
	"goblog/model"
)

type termDao struct {
	c *conn
}

// 返回term的dao实例
func (d *Dao) Term() *termDao {
	if d.term == nil {
		d.term = &termDao{
			c: d.c,
		}
	}
	return d.term
}

// 通过项ID获取项
func (t *termDao) GetByID(id uint32) (*model.Term, bool) {
	term, has := model.Term{}, true
	err := t.c.Get(&term, "SELECT * FROM terms WHERE tid = ?", id)
	if err == sql.ErrNoRows {
		has = false
	} else {
		t.c.panicExistError(err)
	}
	return &term, has
}

// 通过项Slug获取项
func (t *termDao) GetBySlug(slug string) (*model.Term, bool) {
	term, has := model.Term{}, true
	err := t.c.Get(&term, "SELECT * FROM terms WHERE slug = ?", slug)
	if err == sql.ErrNoRows {
		has = false
	} else {
		t.c.panicExistError(err)
	}
	return &term, has
}

// 获取项类型对应的项列表
func (t *termDao) ListByType(termType model.TermType, desc bool, pi, ps uint32) ([]model.Term, bool) {
	var termList []model.Term
	has := true

	var query string
	if desc {
		query = "SELECT * FROM terms WHERE type = ? ORDER BY tid DESC LIMIT ?, ?"
	} else {
		query = "SELECT * FROM terms WHERE type = ? LIMIT ?, ?"
	}

	err := t.c.Select(&termList, query, termType, (pi-1)*ps, ps)
	t.c.panicExistError(err)
	if len(termList) == 0 {
		has = false
	}
	return termList, has
}

// 通过文章ID获取其对应所有项的列表
func (t *termDao) ListByPostID(id uint32, termType model.TermType) ([]model.Term, bool) {
	var termList []model.Term
	has := true
	err := t.c.Select(
		&termList,
		`SELECT * FROM terms WHERE tid in (
			SELECT tid FROM relationships WHERE pid = ?
		) and type = ?`,
		id, termType,
	)
	t.c.panicExistError(err)

	if len(termList) == 0 {
		has = false
	}
	return termList, has
}

// 新增项
func (t *termDao) Add(term *model.Term) (uint32, bool) {
	res, err := t.c.NamedExec("INSERT INTO terms(name, slug, type, description, count) VALUES(:name, :slug, :type, :description, :count)", term)
	t.c.panicExistError(err)

	id, _ := res.LastInsertId()
	ok := cmpRowsAffected(res, 1)
	return uint32(id), ok
}

// 修改项
func (t *termDao) Modify(term *model.Term) bool {
	res, err := t.c.NamedExec("UPDATE terms SET name=:name, description=:description WHERE tid=:tid", term)
	t.c.panicExistError(err)

	return cmpRowsAffected(res, 1)
}

// 统计项对应文章数
func (t *termDao) CountByType(termType model.TermType) uint32 {
	var count uint32 = 0
	err := t.c.Get(&count, "SELECT COUNT(*) FROM terms WHERE type = ?", termType)
	t.c.panicExistError(err)
	return count
}
