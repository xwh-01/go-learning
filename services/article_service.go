// services/article_service.go
package services

import (
	"context"
	"encoding/json"
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"fmt"
	"github.com/go-redis/redis/v8" // 统一v8版本
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// 上下文统一使用context.Background()（或业务ctx）
var ctx = context.Background()

type ArticleService struct{}

// 1. 修复 GetArticleByID 中的Redis调用
func (s *ArticleService) GetArticleByID(id uint) (*models.Article, error) {
	cacheKey := fmt.Sprintf("article:detail:%d", id)

	// 修复：v8的Get需要 (ctx, key) → 参数数量正确
	cachedData, err := global.RedisDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var article models.Article
		if json.Unmarshal([]byte(cachedData), &article) == nil {
			return &article, nil
		}
	} else if err != redis.Nil {
		fmt.Printf("[WARN] Redis查询出错(key:%s): %v\n", cacheKey, err)
	}

	var article models.Article
	if err := global.Db.First(&article, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 修复：v8的Set需要 (ctx, key, value, expire) → 参数数量正确
			global.RedisDB.Set(ctx, cacheKey, "", 3*time.Minute)
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	articleJSON, err := json.Marshal(article)
	if err == nil {
		expireTime := 30*time.Minute + time.Duration(rand.Intn(300))*time.Second
		// 修复：Set参数顺序（ctx在前）
		if err := global.RedisDB.Set(ctx, cacheKey, articleJSON, expireTime).Err(); err != nil {
			fmt.Printf("[WARN] 回种缓存失败(key:%s): %v\n", cacheKey, err)
		}
	}

	return &article, nil
}

// 2. 修复 CreateArticle 中的Del调用
func (s *ArticleService) CreateArticle(article *models.Article) error {
	tx := global.Db.Begin()

	if err := tx.Create(article).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 修复：v8的Del需要 (ctx, keys...) → 先传ctx，再传key
	if err := global.RedisDB.Del(ctx, "articles:list:*").Err(); err != nil {
		fmt.Printf("[WARN] 清理列表缓存失败: %v\n", err)
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// GetArticlesWithPagination 无Redis调用，无需修改
func (s *ArticleService) GetArticlesWithPagination(page, pageSize int) ([]models.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	var articles []models.Article
	var total int64

	if err := global.Db.Model(&models.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := global.Db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
