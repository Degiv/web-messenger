package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"messenger/internals/handlers"
	"messenger/internals/services"
	storage2 "messenger/internals/storage"
)

func NewLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.DebugLevel)
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLogger, _ := config.Build()
	return zapLogger
}

func main() {
	log := NewLogger()
	defer log.Sync() // flushes buffer, if any
	//logSugar := log.Sugar()

	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Error("Failed connect to database", zap.Error(err))
	}

	users := storage2.NewUsers(db)
	messages := storage2.NewMessages(db)
	conferences := storage2.NewConferences(db)

	messengerService := services.NewMessenger(users, messages, conferences)

	e := echo.New()
	e.Use(middleware.Logger())
	handler := handlers.NewMessengerHandler(messengerService, log)
	handler.RegisterRoutes(e)

	err = e.Start(":3000")
	log.Error("Failed listen and serve", zap.Error(err))
}
