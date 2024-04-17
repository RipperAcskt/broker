package handlers

import "github.com/gin-gonic/gin"

type Handlers struct {
	messages *Messages
	exchange *Exchange
}

func (h Handlers) InitRouters() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	h.messages.InitRouters(router)
	h.exchange.InitRouters(router)

	return router
}

func New(broker brokerUsecases) *Handlers {
	return &Handlers{
		messages: NewMessages(broker),
		exchange: NewExchange(broker),
	}
}
