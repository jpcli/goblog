package repository

import (
	"strings"
)

// InsertRelationships insert one or more new relationships into 'relationships' table.
func (m *mysql) InsertRelationships(src ...*Relationship) {
	queryTemp := strings.TrimSuffix(strings.Repeat("(?, ?),", len(src)), ",")
	query := `INSERT INTO relationships(pid, tid) VALUES` + queryTemp

	args := make([]interface{}, 0, len(src)*2)
	for _, v := range src {
		args = append(args, v.Pid, v.Tid)
	}

	m.Exec(query, args...)
}

// DeleteRelationshipsByPid delete one post's all relationships from 'relationships' table.
func (m *mysql) DeleteRelationshipsByPid(pid uint32) {
	query := `DELETE FROM relationships WHERE pid = ?`
	m.Exec(query, pid)
}
