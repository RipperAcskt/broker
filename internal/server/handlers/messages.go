package handlers

import (
	"context"
	"github.com/RipperAcskt/broker/internal/entities/models"
	"github.com/RipperAcskt/broker/internal/server/handlers/dto"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type brokerUsecases interface {
	NewExchange(ctx context.Context, key string) error
	NewQueue(ctx context.Context, key string, offset int) (chan models.Message, int, error)
	NewMessage(ctx context.Context, key string, body any) error
	ReceiveMessage(ctx context.Context, queue chan models.Message) (models.Message, error)
	Disconnect(key string, queueIndex int)
}

type Messages struct {
	upgrader websocket.Upgrader

	brokerUsecases brokerUsecases
}

func (m *Messages) Send(ctx *gin.Context) {
	var message dto.Message

	if err := ctx.BindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := m.brokerUsecases.NewMessage(ctx, message.Key, message.Message); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (m *Messages) Receive(ctx *gin.Context) {
	key := ctx.Param("key")

	var offset int
	var err error
	offsetStr := ctx.Query("offset")
	if offsetStr == "" {
		offset = -1
	} else {
		offset, err = strconv.Atoi(ctx.Query("offset"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	queue, index, err := m.brokerUsecases.NewQueue(ctx, key, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	conn, err := m.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer conn.Close()

	isRunning := true
	conn.SetCloseHandler(func(code int, text string) error {
		isRunning = false
		m.brokerUsecases.Disconnect(key, index)
		return nil
	})

	go func(key string, index int) {
		conn.ReadMessage()
	}(key, index)

	for isRunning {
		message, err := m.brokerUsecases.ReceiveMessage(ctx, queue)
		if err != nil {
			log.Println(err)
		}

		err = conn.WriteJSON(message)
		if err != nil {
			log.Println(err)
		}
	}
}

func (m *Messages) InitRouters(router *gin.Engine) {
	messages := router.Group("/messages")

	messages.POST("/send", m.Send)
	messages.GET("/receive/:key", m.Receive)
}

func NewMessages(brokerUsecases brokerUsecases) *Messages {
	return &Messages{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		brokerUsecases: brokerUsecases,
	}
}
