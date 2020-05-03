package mysql

import (
	"bluebell_backend/logger"
	"bluebell_backend/models"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	sqlStr := `insert into post(
	post_id, caption, content, author_id, community_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, post.PostID, post.Caption,
		post.Content, post.AuthorId, post.CommunityID)
	if err != nil {
		logger.Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

//GetPostByID
func GetPostByID(idStr string) (post *models.ApiPostDetail, err error) {
	post = new(models.ApiPostDetail)
	sqlStr := `select post_id, caption, content, author_id, community_id, create_time
	from post
	where post_id = ?`
	err = db.Get(post, sqlStr, idStr)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		logger.Error("query post failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

func GetPostListByIDs(ids []string) (posts []*models.Post, err error) {
	sqlStr := `select post_id, caption, content, author_id, community_id, create_time
	from post
	where post_id in (?)`
	// 动态填充id
	query, args, err := sqlx.In(sqlStr, ids)
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&posts, query, args...)
	return
}
