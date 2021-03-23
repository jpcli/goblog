package repository

import (
	"database/sql"
	"fmt"
	"goblog/utils/errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

// db is global database connection
var db *sqlx.DB

func OpenDatabase(ip, port, usr, pwd, database string) {
	d, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
		usr, pwd, ip, port, database,
	))
	if err != nil {
		panic(errors.NewErrorWithStack("failed to open database"))
	}
	db = d
}

type Repository struct {
	*mysql
}

func NewRepository() *Repository {
	return &Repository{
		&mysql{
			db: db,
		},
	}
}

// mysql struct is used to access tha mysql database.
type mysql struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	err error
}

// GetError returns the err of mysql.
func (m *mysql) GetError() error {
	return m.err
}

// Begin a new transaction.
func (m *mysql) Begin() {
	m.tx, m.err = m.db.Beginx()
	m.err = errors.WrapfErrorWithStack(m.err, "failed to begin a transaction")
}

// Commit the transaction which has been began before.
func (m *mysql) Commit() {
	if m.tx == nil {
		m.err = errors.NewErrorWithStack("there is no tx")
	}

	// empty the TX
	defer func() {
		m.tx = nil
	}()

	var err error
	if m.err != nil {
		err = m.tx.Rollback()
		if err != nil {
			m.err = errors.AttachErrorMessage(m.err, "failed to rollback: %s", err.Error())
		}
		return
	}

	err = m.tx.Commit()
	if err != nil {
		m.err = errors.WrapfErrorWithStack(err, "failed to commit")
	}
}

// Rollback used to rollback tx manually and set tx nil.
func (m *mysql) Rollback() error {
	if m.tx == nil {
		return nil
	}

	// empty the TX
	defer func() {
		m.tx = nil
	}()

	err := m.tx.Rollback()
	return errors.WrapfErrorWithStack(err, "failed to rollback")
}

// Get execute a query and scan a row into dest.
// Error will happen when no row has selected.
func (m *mysql) Get(dest interface{}, query string, args ...interface{}) {
	if m.err != nil {
		return
	}

	var err error
	if m.tx != nil {
		err = m.tx.Get(dest, query, args...)
	} else {
		err = m.db.Get(dest, query, args...)
	}

	if err == sql.ErrNoRows {
		m.err = errors.WrapfErrorWithStack(ErrNoRow, "failed to get")
	} else {
		m.err = errors.WrapfErrorWithStack(err, "failed to get")
	}
}

// Select execute a query and scan each row into dest.
func (m *mysql) Select(dest interface{}, query string, args ...interface{}) {
	if m.err != nil {
		return
	}

	var err error
	if m.tx != nil {
		err = m.tx.Select(dest, query, args...)
	} else {
		err = m.db.Select(dest, query, args...)
	}
	m.err = errors.WrapfErrorWithStack(err, "failed to select")
}

// Exec execute a query that doesn't return rows
func (m *mysql) Exec(query string, args ...interface{}) (rowsAffected, lastInsertID int64) {
	if m.err != nil {
		return
	}

	var result sql.Result
	var err error
	if m.tx != nil {
		result, err = m.tx.Exec(query, args...)
	} else {
		result, err = m.db.Exec(query, args...)
	}
	if err != nil {
		m.err = errors.WrapfErrorWithStack(err, "failed to execute")
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		m.err = errors.WrapfErrorWithStack(err, "failed to get rowsAffected")
		return
	} else if rowsAffected == 0 {
		m.err = errors.WrapErrorWithStack(ErrNoRowAffected)
		return
	}

	if strings.HasPrefix(strings.ToLower(query), "insert") {
		if lastInsertID, err = result.LastInsertId(); err != nil {
			m.err = errors.WrapfErrorWithStack(err, "failed to get lastInsertId")
		}
	}
	return
}

var (
	ErrNoRow           = fmt.Errorf("no row be selected")
	ErrNoRowAffected   = fmt.Errorf("no row has been affected")
	ErrUnexpectedValue = fmt.Errorf("receive unexpected value")
)
