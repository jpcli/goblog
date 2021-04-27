package dao

import "goblog/model"

type relaDao struct {
	c *conn
}

// 返回relationship的dao实例
func (d *Dao) Rela() *relaDao {
	if d.rela == nil {
		d.rela = &relaDao{
			c: d.c,
		}
	}
	return d.rela
}

// 新增一条关系
func (r *relaDao) Add(rela *model.Relationship) bool {
	res, err := r.c.NamedExec("INSERT INTO relationships(pid, tid) VALUES(:pid, :tid)", rela)
	panicExistError(err)
	return cmpRowsAffected(res, 1)
}

// 删除文章ID对应所有关系
func (r *relaDao) RemoveByPostID(id uint32) bool {
	res, err := r.c.Exec("DELETE FROM relationships WHERE pid = ?", id)
	panicExistError(err)
	affected, _ := res.RowsAffected()
	return affected > 0
}
