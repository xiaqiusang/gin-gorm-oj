package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatPost(p *models.Post) (err error) {
	//生成请求的post id
	p.ID = int64(snowflake.GenID())
	//保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatPost(p.ID)
	return
	//返回
}

func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	//根据作者id拿到作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(pid) failed", zap.Int64("post.AuthorID", post.AuthorID), zap.Error(err))
		return
	}
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(pid) failed", zap.Int64("post.CommunityID", post.CommunityID), zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}

	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者id拿到作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(pid) failed", zap.Int64("post.AuthorID", post.AuthorID), zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(pid) failed", zap.Int64("post.CommunityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn(" redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	//根据ID去数据库里查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	//将帖子的作者及分区信息查询出来填充到帖子
	for idx, post := range posts {
		//根据作者id拿到作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(pid) failed", zap.Int64("post.AuthorID", post.AuthorID), zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(pid) failed", zap.Int64("post.CommunityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn(" redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	//根据ID去数据库里查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	//将帖子的作者及分区信息查询出来填充到帖子
	for idx, post := range posts {
		//根据作者id拿到作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(pid) failed", zap.Int64("post.AuthorID", post.AuthorID), zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(pid) failed", zap.Int64("post.CommunityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// 将两个查询接口合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//根据请求参数不同，执行不同的逻辑
	if p.CommunityID == 0 {
		//查所有
		data, err = GetPostList2(p)
	} else {
		//根据社区id查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
	}
	return
}
