// controllers/like_controller.go
package controllers

import (
	"exchangeapp/services"
	"exchangeapp/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func LikeArticle(ctx *gin.Context) {
	// 1. 获取并转换文章ID
	idStr := ctx.Param("id")
	articleID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "无效的文章ID")
		return
	}

	// 2. 调用Service处理核心业务
	likeService := new(services.LikeService)
	if err := likeService.LikeArticle(uint(articleID)); err != nil {
		utils.InternalServerError(ctx, "点赞失败: "+err.Error())
		return
	}

	// 3. 统一成功响应
	utils.SuccessWithMessage(ctx, "点赞成功")
}

func GetArticlesLikes(ctx *gin.Context) {
	idStr := ctx.Param("id")
	articleID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "无效的文章ID")
		return
	}

	likeService := new(services.LikeService)
	likes, err := likeService.GetArticleLikes(uint(articleID))
	if err != nil {
		utils.InternalServerError(ctx, "获取点赞数失败")
		return
	}

	utils.Success(ctx, gin.H{"likes": likes})
}
