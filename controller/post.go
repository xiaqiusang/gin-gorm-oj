package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

const (
	orderTime  = "time"
	orderScore = "score"
)

// CreatePostHandle 创建帖子的函数
func CreatePostHandle(c *gin.Context) {
	//1.获取参数及参数的校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error ", zap.Any("err", err))
		zap.L().Error("Create Post with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从c取到当前发请求的用户ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLlogin)
		return
	}
	p.AuthorID = userID
	//2.创建帖子
	if err := logic.CreatPost(p); err != nil {
		zap.L().Error("CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context) {
	//1.获取参数(URL中获取帖子的ID)
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据id取出帖子的数据(查数据库)
	date, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("get post detail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, date)
}

// GetPostListHandler 获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	page, size := getPageInfo(c)
	//获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表接口
// 根据前端传来的参数动态的获取帖子列表
// 按创建时间排序 或  者按照分数排序
// 1.获取参数
// 2.去redis查询id列表
// 3.根据id去数据库查询帖子详细信息
func GetPostListHandler2(c *gin.Context) {
	//Get请求参数:/api/v1/posts2?page=1&size=10&order=time
	P := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	//获取分页参数
	if err := c.ShouldBindQuery(P); err != nil {
		zap.L().Error(" GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//c.ShouldBind()根据请求的数据类型选择响应的方法去获取数据
	//c.ShouldBindJSON()如果请求中携带的是json格式的数据，才能使这个方法获取到数据

	//获取数据
	data, err := logic.GetPostList2(P)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

//// 根据社区去查询帖子列表
//func GetCommentPostListHandler(c *gin.Context) {
//	//Get请求参数:/api/v1/posts2?page=1&size=10&order=time
//	P := &models.ParamCommunityPostList{
//		Page:  1,
//		Size:  10,
//		Order: models.OrderTime,
//	}
//	//获取分页参数
//	if err := c.ShouldBindQuery(P); err != nil {
//		zap.L().Error(" GetCommentPostList with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//	//c.ShouldBind()根据请求的数据类型选择响应的方法去获取数据
//	//c.ShouldBindJSON()如果请求中携带的是json格式的数据，才能使这个方法获取到数据
//	//获取数据
//	data, err := logic.GetCommentPostList2(p * models.ParamPostList{})
//	if err != nil {
//		zap.L().Error("logic.GetPostList failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	//返回响应
//	ResponseSuccess(c, data)
//}
