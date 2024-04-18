package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/RipperAcskt/broker/internal/server/handlers"
	"github.com/RipperAcskt/broker/internal/usecases"
	"go.uber.org/zap"
	"net/http"

	"github.com/RipperAcskt/broker/internal/repository/mongo"
	"github.com/RipperAcskt/broker/internal/server"
	"github.com/RipperAcskt/broker/internal/service"
)

type App struct {
	errChan chan error
}

func New() *App {
	return &App{
		errChan: make(chan error),
	}
}

func (a *App) Run() error {
	client, err := mongo.New()
	if err != nil {
		return fmt.Errorf("could not create mongo client: %w", err)
	}

	exchangeRepo := mongo.NewExchange(client)
	queueRepo := mongo.NewQueue(client)
	messageRepo := mongo.NewMessage(client)

	exchangeService, err := service.NewExchange(exchangeRepo)
	if err != nil {
		return fmt.Errorf("could not create exchange service: %w", err)
	}
	queueService, err := service.NewQueue(queueRepo)
	if err != nil {
		return fmt.Errorf("could not create queue service: %w", err)
	}
	messageService := service.NewMessages(messageRepo)

	broker, err := usecases.NewBroker(context.TODO(), exchangeService, queueService, messageService)

	handl := handlers.New(broker)

	srv := server.New()
	srv.Log, _ = zap.NewProduction()

	go func(errChan chan error) {
		if err := srv.Run(handl.InitRouters()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("server run failed: %w", err)
		}
	}(a.errChan)

	go func(errChan chan error) {
		if err := srv.WaitForShutDown(); err != nil {
			errChan <- fmt.Errorf("shut down failed: %w", err)
		}
		close(errChan)
	}(a.errChan)

	for err := range a.errChan {
		return err
	}
	return nil
}
