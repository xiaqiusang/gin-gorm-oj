package controller

import (
	"bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

//---跟社区相关的

func CommunityHandler(c *gin.Context) {
	//查询到所有的社区(community_id,community_name)以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端暴露给外面
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	//获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64) //10进制，64位，校验作用查看输入类型是否有效
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//查询到所有的社区(community_id,community_name)以列表的形式返回
	data, err := logic.GetCommunityDetail(id) //data为返回id搜索到的详细信息
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端暴露给外面
		return
	}
	ResponseSuccess(c, data)
}
