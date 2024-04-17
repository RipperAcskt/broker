package handlers

import (
	"github.com/RipperAcskt/broker/internal/server/handlers/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Exchange struct {
	brokerUsecases brokerUsecases
}

func (e *Exchange) Create(ctx *gin.Context) {
	var exchange dto.Exchange

	if err := ctx.BindJSON(&exchange); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := e.brokerUsecases.NewExchange(ctx, exchange.Key); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (e *Exchange) InitRouters(router *gin.Engine) {
	messages := router.Group("/exchanges")

	messages.POST("/", e.Create)
}

func NewExchange(brokerUsecases brokerUsecases) *Exchange {
	return &Exchange{
		brokerUsecases: brokerUsecases,
	}
}
