package repository

import (
	"github.com/spf13/cast"
)

// InsertAttachment insert a new attachment to 'attachments' table.
func (m *mysql) InsertAttachment(src *Attachment) (aid uint32) {
	query := `INSERT INTO attachments(created, relativePath) VALUES(?, ?)`
	_, id := m.Exec(query, src.Created, src.RelativePath)
	return cast.ToUint32(id)
}
