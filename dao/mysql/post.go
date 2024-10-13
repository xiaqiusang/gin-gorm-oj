package mysql

import (
	"bluebell/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id,title,content,author_id,community_id)values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostById 根据id查询单个帖子数据
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id,title,content,author_id ,community_id,create_time from post where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post ORDER BY post.create_time DESC limit ?,?`
	posts = make([]*models.Post, 0, 2) //长度为0，容量为2 ，如果写成([]*models.Post, 2)就是长度和容量都为2
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 根据id查询帖子
func GetPostListByIDs(ids []string) (postsList []*models.Post, err error) {
	sqlStr := `select post_id ,title,content,author_id,community_id,create_time from post where post_id in (?) order by FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ",")) //先转换sql语句
	//query的内容为select post_id ,title,content,author_id,community_id,create_time from post where id in (?,?,?)，为ids的数量
	//args的内容为[]interface{}{"1", "2", "3"}，为ids的数量
	if err != nil {
		return nil, err
	}
	//这段代码用于处理占位符格式的转换，以确保生成的 SQL 查询能够在特定的数据库驱动下正确执行。
	query = db.Rebind(query)
	err = db.Select(&postsList, query, args...)
	return
}
