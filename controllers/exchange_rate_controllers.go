package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.Exchangerate

	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		utils.HandleAppError(ctx, utils.ErrInvalidParams)
		return
	}

	exchangeRate.Date = time.Now()

	if err := global.Db.Create(&exchangeRate).Error; err != nil {
		utils.HandleAppError(ctx, utils.ErrInternalServer)
		return
	}

	utils.Success(ctx, exchangeRate)
}

func GetexchangeRates(ctx *gin.Context) {
	var exchangeRates []models.Exchangerate

	if err := global.Db.Find(&exchangeRates).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleAppError(ctx, utils.ErrRecordNotFound)
		} else {
			utils.HandleAppError(ctx, utils.ErrInternalServer)
		}
		return
	}

	utils.Success(ctx, exchangeRates)
}
