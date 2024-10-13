package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 投票
func PostVoteController(c *gin.Context) {
	//参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(&p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) //翻译并去掉结构体中的错误标识
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLlogin)
		return
	}
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
