package dao

import "goblog/model"

type attaDao struct {
	c *conn
}

// 返回attachment的dao实例
func (d *Dao) Atta() *attaDao {
	if d.atta == nil {
		d.atta = &attaDao{
			c: d.c,
		}
	}
	return d.atta
}

// 新增附件
func (a *attaDao) Add(atta *model.Attachment) (uint32, bool) {
	res, err := a.c.NamedExec("INSERT INTO attachments(created, relativePath) VALUES(:created, :relativePath)", atta)
	panicExistError(err)
	id, _ := res.LastInsertId()
	ok := cmpRowsAffected(res, 1)
	return uint32(id), ok
}
