package repository

import (
	"fmt"

	"github.com/spf13/cast"
)

// InsertPost inserts a new post into 'posts' table.
func (m *mysql) InsertPost(src *Post) (pid uint32) {
	// if s := src.Status; s != PostStatusPublish && s != PostStatusSticky {
	// 	err := m.Rollback()
	// 	if err != nil {
	// 		panic(err)
	// 	} else {
	// 		panic(errors.NewErrorWithStack("unexpected post status"))
	// 	}
	// }

	query := `INSERT INTO posts(title, created, modified, keywords, excerpt, text, status, commentCount, commentAllow)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, id := m.Exec(query, src.Title, src.Created, src.Modified, src.Keywords, src.Excerpt, src.Text, src.Status, src.CommentCount, src.CommentAllow)
	return cast.ToUint32(id)
}

// UpdatePost updates title, modified, description, keywords, text and commentAllow of post.
func (m *mysql) UpdatePost(pid uint32, src *Post) {
	query := `UPDATE posts
	SET title = ?, modified = ?, excerpt = ?, keywords = ?, text = ?, commentAllow = ?
	WHERE pid = ?`
	m.Exec(query, src.Title, src.Modified, src.Excerpt, src.Keywords, src.Text, src.CommentAllow, pid)
}

// UpdatePostStatus updates status of post.
func (m *mysql) UpdatePostStatus(pid uint32, status string) {
	query := `UPDATE posts SET status = ? WHERE pid = ?`
	m.Exec(query, status, pid)
}

// IncreasePostCommentCountByPid increse commentCount of post by 1.
// func (m *mysql) IncreasePostCommentCountByPid(pid uint32) {
// 	query := `UPDATE posts SET commentCount = commentCount + 1 WHERE pid = ?`
// 	m.Exec(query, pid)
// }

func (m *mysql) GetPostByPid(fields string, pid uint32) *Post {
	var ret Post
	query := fmt.Sprintf("SELECT %s FROM posts WHERE pid = ?", fields)
	m.Get(&ret, query, pid)
	return &ret
}

func (m *mysql) GetPostListByStatus(fields, status string, offset, count uint32) []Post {
	var ret []Post
	query := fmt.Sprintf("SELECT %s FROM posts WHERE status = ? ORDER BY pid DESC LIMIT ?, ?", fields)
	m.Select(&ret, query, status, offset, count)
	return ret
}

func (m *mysql) GetPostList(fields string, offset, count uint32) []Post {
	var ret []Post
	query := fmt.Sprintf("SELECT %s FROM posts WHERE status != 'trash' ORDER BY pid DESC LIMIT ?, ?", fields)
	m.Select(&ret, query, offset, count)
	return ret
}

func (m *mysql) GetPostListByTid(fields string, tid uint32, offset, count uint32) []Post {
	var ret []Post
	query := fmt.Sprintf(`SELECT %s FROM posts 
	WHERE pid in (
		SELECT pid FROM relationships WHERE tid = ?
	) and status != 'trash' 
	ORDER BY pid DESC LIMIT ?, ?`, fields)
	m.Select(&ret, query, tid, offset, count)
	return ret
}

func (m *mysql) CountPost(status string) (count uint32) {
	query := "SELECT COUNT(*) FROM posts WHERE status = ?"
	m.Get(&count, query, status)
	return
}

func (m *mysql) GetLastModified() (modified uint32) {
	query := "SELECT modified FROM posts ORDER BY modified LIMIT 1"
	m.Get(&modified, query)
	return
}

// InsertPostmeta inserts a new postmeta into 'postmeta' table.
func (m *mysql) InsertPostmeta(src *Postmeta) (mid uint32) {
	query := `INSERT INTO postmeta(pid, metaKey, metaValue) VALUES(?, ?, ?)`
	_, id := m.Exec(query, src.Pid, src.MetaKey, src.MetaValue)
	return cast.ToUint32(id)
}

// UpdatePostmetaValue updates the metaValue of postmeta in 'postmeta' table.
func (m *mysql) UpdatePostmetaValue(pid uint32, metaKey, metaValue string) {
	query := `UPDATE postmeta SET metaValue = ? WHERE pid = ? and metaKey = ?`
	m.Exec(query, metaValue, pid, metaKey)
}

func (m *mysql) GetPostmetaValue(pid uint32, metaKey string) (metaValue string) {
	query := `SELECT metaValue FROM postmeta WHERE pid = ? and metaKey = ?`
	m.Get(&metaValue, query, pid, metaKey)
	return
}

func (m *mysql) IncrPostViewCountByPid(pid uint32) {
	query := `UPDATE postmeta SET metaValue = metaValue + 1 WHERE pid = ? and metaKey = 'post_view_count'`
	m.Exec(query, pid)
}
