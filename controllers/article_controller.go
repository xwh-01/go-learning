package controllers

import (
	"errors"
	"exchangeapp/models"
	"exchangeapp/services"
	"exchangeapp/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func CreateArticle(ctx *gin.Context) {
	var article models.Article
	if err := ctx.ShouldBindJSON(&article); err != nil {
		utils.HandleAppError(ctx, utils.ErrInvalidParams)
		return
	}

	articleService := new(services.ArticleService)
	if err := articleService.CreateArticle(&article); err != nil {
		utils.HandleAppError(ctx, utils.ErrInternalServer)
		return
	}

	utils.Success(ctx, article)
}

func GetArticlesByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.HandleAppError(ctx, utils.ErrInvalidParams)
		return
	}

	articleService := new(services.ArticleService)
	article, err := articleService.GetArticleByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "文章不存在" {
			utils.HandleAppError(ctx, utils.ErrRecordNotFound)
		} else {
			utils.HandleAppError(ctx, utils.ErrInternalServer)
		}
		return
	}

	utils.Success(ctx, article)
}

func GetArticles(ctx *gin.Context) {
	var query models.ArticleQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.HandleAppError(ctx, utils.ErrInvalidParams)
		return
	}
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	articleService := new(services.ArticleService)
	articles, total, err := articleService.GetArticlesWithPagination(query.Page, query.PageSize)
	if err != nil {
		utils.HandleAppError(ctx, utils.ErrInternalServer)
		return
	}

	utils.Success(ctx, gin.H{
		"items": articles,
		"pagination": gin.H{
			"page":       query.Page,
			"page_size":  query.PageSize,
			"total":      total,
			"total_page": (int(total) + query.PageSize - 1) / query.PageSize,
		},
	})
}
