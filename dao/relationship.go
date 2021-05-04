package dao

import (
	"fmt"
	"strings"
)

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

// 新增多条关系
func (r *relaDao) AddMore(pid uint32, tids ...uint32) bool {
	// 插入n条关系记录
	var s strings.Builder
	for _, tid := range tids {
		res, err := r.c.Exec("INSERT INTO relationships(pid, tid) VALUES(?, ?)", pid, tid)
		r.c.panicExistError(err)
		if !cmpRowsAffected(res, 1) {
			return false
		}
		fmt.Fprintf(&s, "%d,", tid)
	}

	// 对应项的count+1
	query := fmt.Sprintf("UPDATE terms SET count=count+1 WHERE tid in (%s)", strings.TrimSuffix(s.String(), ","))
	res, err := r.c.Exec(query)
	r.c.panicExistError(err)
	return cmpRowsAffected(res, int64(len(tids)))
}

// 删除文章ID对应所有关系
func (r *relaDao) RemoveByPostID(id uint32) bool {
	// 对应文章对应项的count-1
	res, err := r.c.Exec(`UPDATE terms SET count=count-1 WHERE tid in (SELECT tid FROM relationships WHERE pid = ?)`, id)
	r.c.panicExistError(err)
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return false
	}

	// 删除文章-项关系
	res, err = r.c.Exec("DELETE FROM relationships WHERE pid = ?", id)
	r.c.panicExistError(err)
	affected, _ = res.RowsAffected()
	return affected > 0
}
