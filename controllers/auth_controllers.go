package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.HandleAppError(ctx, utils.ErrInvalidParams)
		return
	}

	hashedPwd, err := utils.HashPassword(user.Password)

	if err != nil {
		utils.HandleAppError(ctx, utils.ErrInternalServer)
		return
	}
	user.Password = hashedPwd

	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		utils.HandleAppError(ctx, utils.ErrInternalServer)
		return
	}

	if err := global.Db.Create(&user).Error; err != nil {
		utils.HandleAppError(ctx, utils.ErrUserExists)
		return
	}

	utils.Success(ctx, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.HandleAppError(ctx, utils.ErrInvalidParams)
		return
	}

	var user models.User
	if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		utils.HandleAppError(ctx, utils.ErrInvalidPassword)
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		utils.HandleAppError(ctx, utils.ErrInvalidPassword)
		return
	}

	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		utils.HandleAppError(ctx, utils.ErrInternalServer)
		return
	}

	utils.Success(ctx, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}
