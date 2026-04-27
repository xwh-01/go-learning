// services/like_service.go
package services

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"fmt"
	"time"
)

type LikeService struct{}

// LikeArticle 点赞文章，保证Redis和MySQL数据同步
func (s *LikeService) LikeArticle(articleID uint) error {
	// 1. 验证文章是否存在
	var article models.Article
	if err := global.Db.First(&article, articleID).Error; err != nil {
		return errors.New("文章不存在")
	}

	// 2. 开启数据库事务
	tx := global.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 3. 【核心】在事务内更新MySQL点赞数
	if err := article.IncrementLikes(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("更新数据库点赞数失败: %w", err)
	}

	// 4. 更新Redis点赞数
	redisKey := fmt.Sprintf("article:%d:likes", articleID)
	if err := global.RedisDB.Incr(ctx, redisKey).Err(); err != nil {
		tx.Rollback() // Redis失败，回滚MySQL事务
		return fmt.Errorf("更新Redis点赞数失败: %w", err)
	}

	// 5. 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// GetArticleLikes 获取文章点赞数（优先从Redis读，没有则读MySQL并回种）
func (s *LikeService) GetArticleLikes(articleID uint) (int64, error) {
	redisKey := fmt.Sprintf("article:%d:likes", articleID)

	// 先查Redis
	likesStr, err := global.RedisDB.Get(ctx, redisKey).Result()
	if err == nil {
		// Redis存在，解析返回
		var likes int64
		fmt.Sscan(likesStr, &likes)
		return likes, nil
	}

	// Redis不存在（过期或首次），查数据库
	var article models.Article
	if err := global.Db.First(&article, articleID).Error; err != nil {
		return 0, errors.New("文章不存在")
	}

	// 将数据库中的值回种到Redis，设置一定过期时间（如1小时）
	if err := global.RedisDB.Set(ctx, redisKey, article.Likes, time.Hour).Err(); err != nil {
		// 回种失败只记录日志，不影响主流程
		fmt.Printf("回种Redis失败: %v\n", err)
	}

	return int64(article.Likes), nil
}
