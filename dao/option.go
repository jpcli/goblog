package dao

import (
	"database/sql"
	"goblog/model"
)

type optDao struct {
	c *conn
}

// 返回option的dao实例
func (d *Dao) Opt() *optDao {
	if d.opt == nil {
		d.opt = &optDao{
			c: d.c,
		}
	}
	return d.opt
}

// 获取name对应选项的值
func (o *optDao) GetByName(name string) (string, bool) {
	value, has := "", true
	err := o.c.Get(&value, "SELECT value FROM options WHERE name = ?", name)
	if err == sql.ErrNoRows {
		has = false
	} else {
		panicExistError(err)
	}
	return value, has
}

// 修改选项
func (o *optDao) Modify(opt *model.Option) bool {
	res, err := o.c.NamedExec("UPDATE options SET value=:value WHERE name=:name", opt)
	panicExistError(err)
	return cmpRowsAffected(res, 1)
}
