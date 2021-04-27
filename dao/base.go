package dao

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// 数据库连接池单例
var db *sqlx.DB

func OpenDatabase(ip string, port int, usr, pwd, database string) {
	d, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4", usr, pwd, ip, port, database))
	if err != nil {
		panic(fmt.Errorf("打开数据库失败"))
	}
	db = d
}

// conn用于控制使用db或tx，一次dao维护一个conn
type conn struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

type Dao struct {
	c    *conn
	post *postDao
	term *termDao
	rela *relaDao
	atta *attaDao
	opt  *optDao
}

func NewDao() *Dao {
	return &Dao{
		c: &conn{
			db: db,
		},
	}
}

// 为conn封装sqlx操作，用于控制执行的连接是db还是tx
// ------conn方法定义开始------

func (c *conn) Get(dest interface{}, query string, args ...interface{}) error {
	if c.tx != nil {
		return c.tx.Get(dest, query, args...)
	} else {
		return c.db.Get(dest, query, args...)
	}
}

func (c *conn) Select(dest interface{}, query string, args ...interface{}) error {
	if c.tx != nil {
		return c.tx.Select(dest, query, args...)
	} else {
		return c.db.Select(dest, query, args...)
	}
}

func (c *conn) Exec(query string, args ...interface{}) (sql.Result, error) {
	if c.tx != nil {
		return c.tx.Exec(query, args...)
	} else {
		return c.db.Exec(query, args...)
	}
}

func (c *conn) NamedExec(query string, arg interface{}) (sql.Result, error) {
	if c.tx != nil {
		return c.tx.NamedExec(query, arg)
	} else {
		return c.db.NamedExec(query, arg)
	}
}

func (c *conn) Begin() error {
	var err error
	if c.tx == nil {
		c.tx, err = c.db.Beginx()
	} else {
		err = fmt.Errorf("上一个事务未提交或回滚，不能再开始新事务")
	}
	return err
}

func (c *conn) Commit() error {
	if c.tx == nil {
		return fmt.Errorf("没有事务需要提交")
	}

	err := c.tx.Commit()
	// 如果commit成功，就删除当前事务指针；否则留待rollback进行删除
	if err == nil {
		c.tx = nil
	}
	return err
}

func (c *conn) Rollback() error {
	if c.tx == nil {
		return fmt.Errorf("没有事务需要回退")
	}

	err := c.tx.Rollback()
	c.tx = nil
	return err
}

// ------conn方法定义结束------

// 比较影响的行数，相同返回真，不相同返回假
func cmpRowsAffected(res sql.Result, expected int64) bool {
	ok := true
	if affected, _ := res.RowsAffected(); affected != expected {
		ok = false
	}
	return ok
}

// 如果Error不为nil，则panic
func panicExistError(err error) {
	if err != nil {
		panic(err)
	}
}
