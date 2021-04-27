package dao

import (
	"database/sql"
	"fmt"
	"goblog/model"
)

type postDao struct {
	c        *conn
	postmeta *postmetaDao
}
type postmetaDao struct {
	c *conn
}

// 返回post的dao实例
func (d *Dao) Post() *postDao {
	if d.post == nil {
		d.post = &postDao{
			c: d.c,
		}
	}
	return d.post
}

// 返回postmeta的dao实例
func (p *postDao) Meta() *postmetaDao {
	if p.postmeta == nil {
		p.postmeta = &postmetaDao{
			c: p.c,
		}
	}
	return p.postmeta
}

// 获取一篇文章
func (p *postDao) GetByID(id uint32) (*model.Post, bool) {
	post, has := model.Post{}, true
	err := p.c.Get(&post, "SELECT * FROM posts WHERE pid = ?", id)
	if err == sql.ErrNoRows {
		has = false
	} else {
		p.c.panicExistError(err)
	}
	return &post, has
}

// 根据文章类型获取文章列表
func (p *postDao) ListByStatus(status model.PostStatusType, pi, ps uint32) ([]model.Post, bool) {
	var postList []model.Post
	has := true

	err := p.c.Select(&postList, "SELECT * FROM posts WHERE status = ? ORDER BY pid DESC LIMIT ?, ?", status, (pi-1)*ps, ps)
	p.c.panicExistError(err)

	if len(postList) == 0 {
		has = false
	}
	return postList, has
}

// 获取所有文章的归档列表
func (p *postDao) ListArchive() {
	// TODO 归档列表
}

// 根据项ID获取文章列表
func (p *postDao) ListByTermID(id, pi, ps uint32) ([]model.Post, bool) {
	var PostList []model.Post
	has := true
	query := fmt.Sprintf(
		`SELECT * FROM posts 
		WHERE pid in (
			SELECT pid FROM relationships WHERE tid = ?
		) and status in (%d , %d)  
		ORDER BY pid DESC LIMIT ?, ?`,
		model.POST_STATUS_PUBLISH, model.POST_STATUS_STICKY,
	)
	err := p.c.Select(&PostList, query, id, (pi-1)*ps, ps)
	p.c.panicExistError(err)

	if len(PostList) == 0 {
		has = false
	}
	return PostList, has
}

// 新建文章，返回新文章id与是否成功
func (p *postDao) Add(post *model.Post) (uint32, bool) {
	res, err := p.c.NamedExec(
		`INSERT INTO posts(title, created, modified, excerpt, keywords, text, status, commentCount, commentAllow)
		VALUES(:title, :created, :modified, :excerpt, :keywords, :text, :status, :commentCount, :commentAllow)`,
		post,
	)
	p.c.panicExistError(err)

	id, _ := res.LastInsertId()
	ok := cmpRowsAffected(res, 1)
	return uint32(id), ok
}

// 修改文章
func (p *postDao) Modify(post *model.Post) bool {
	res, err := p.c.NamedExec(
		"UPDATE posts SET title=:title, modified=:modified, excerpt=:excerpt, keywords=:keyword, text=:text, commentAllow=:commentAllow WHERE pid=:pid",
		post,
	)
	p.c.panicExistError(err)

	return cmpRowsAffected(res, 1)
}

// 修改文章状态
func (p *postDao) ModifyStatus(id uint32, status model.PostStatusType) bool {
	res, err := p.c.Exec("UPDATE posts SET status = ? WHERE pid = ?", status, id)
	p.c.panicExistError(err)

	return cmpRowsAffected(res, 1)
}

// 统计文章数目
func (p *postDao) CountByStatus(status model.PostStatusType) uint32 {
	var count uint32 = 0
	err := p.c.Get(&count, "SELECT COUNT(*) FROM posts WHERE status = ?", status)
	p.c.panicExistError(err)

	return count
}

// 获取最后修改文章的日期
func (p *postDao) GetLastModified() uint32 {
	var lastModified uint32 = 0
	err := p.c.Get(&lastModified, "SELECT modified FROM posts ORDER BY modified LIMIT 1")
	p.c.panicExistError(err)

	return lastModified
}

// 获取文章元数据
func (m *postmetaDao) GetByKey(id uint32, key string) string {
	var value string
	err := m.c.Get(&value, "SELECT metaValue FROM postmeta WHERE pid = ? and metaKey = ?", id, key)
	if err == sql.ErrNoRows {
		return ""
	} else {
		m.c.panicExistError(err)
	}
	return value
}

// 新增文章元数据
func (m *postmetaDao) Add(id uint32, key, value string) bool {
	res, err := m.c.Exec("INSERT INTO postmeta(pid, metaKey, metaValue) VALUES(?, ?, ?)", id, key, value)
	m.c.panicExistError(err)

	return cmpRowsAffected(res, 1)
}

// 修改文章元数据
func (m *postmetaDao) Modify(id uint32, key, value string) bool {
	res, err := m.c.Exec("UPDATE postmeta SET metaValue = ? WHERE pid = ? and metaKey = ?", value, id, key)
	m.c.panicExistError(err)

	return cmpRowsAffected(res, 1)
}
