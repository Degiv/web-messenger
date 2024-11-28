package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"messenger/internals/handlers"
	"messenger/internals/services/auth"
	"messenger/internals/services/messenger"
	"messenger/internals/storage"
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

	users := storage.NewUsers(db)
	messages := storage.NewMessages(db)
	conferences := storage.NewConferences(db)

	messengerService := messenger.NewMessenger(users, messages, conferences)
	authService := auth.NewAuthService(users)

	e := echo.New()
	e.Use(middleware.Logger())
	handler := handlers.NewMessengerHandler(messengerService, authService, log)
	handler.RegisterRoutes(e)

	err = e.Start(":3000")
	log.Error("Failed listen and serve", zap.Error(err))
}
