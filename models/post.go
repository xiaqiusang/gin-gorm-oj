package models

import "time"

// 内存对齐的概念

type Post struct {
	ID          int64     `json:"id" db:"post_id"`
	AuthorID    int64     `json:"author" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ApiPostDetail 帖子详细信息
type ApiPostDetail struct {
	AuthorName       string `json:"author_name"`
	VoteNum          int64  `json:"votes_num"`
	*Post                   //嵌入帖子结构体
	*CommunityDetail `json:"community_detail"`
}
