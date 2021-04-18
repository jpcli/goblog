package repository

import (
	"goblog/utils/errors"
)

func (m *mysql) SetOption(key, value string) {
	if m.err != nil {
		return
	}
	query := `INSERT INTO options(optionKey, optionValue)VALUES(?, ?) ON DUPLICATE KEY UPDATE optionValue = ?`
	m.Exec(query, key, value, value)
	if errors.Cause(m.err) == ErrNoRowAffected {
		m.err = nil
	}
}

func (m *mysql) GetOption(key string) (value string) {
	query := `SELECT optionValue FROM options WHERE optionKey = ?`
	m.Get(&value, query, key)
	return
}
