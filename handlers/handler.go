package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"messenger/domain"
	"net/http"
	"strconv"
	"time"
)

type MessengerService interface {
	GetConferenceMessages(conferenceID int64) ([]*domain.Message, error)
	GetConferencesByUser(userID int64) ([]*domain.Conference, error)
	CreateConference(usersIDs []int64, name string, createdBy int64, createdAt time.Time) error
}

type MessengerHandler struct {
	service MessengerService
	log     *zap.Logger
}

func NewMessengerHandler(service MessengerService, log *zap.Logger) *MessengerHandler {
	return &MessengerHandler{
		service: service,
		log:     log,
	}
}

func (handler *MessengerHandler) RegisterRoutes(e *echo.Echo) *echo.Echo {
	e.GET("messenger/conferences/{id}", handler.getMessages)
	e.GET("messenger/conferences", handler.getConferences)
	e.POST("messenger/conferences", handler.createConference)
	return e
}

func (handler *MessengerHandler) getMessages(c echo.Context) error {
	conferenceID, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
	if err != nil {
		handler.log.Error("Wrong ID type", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong ID type")
	}

	messages, err := handler.service.GetConferenceMessages(conferenceID)
	if err != nil {
		handler.log.Error("Failed to get messages by conference ID", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to get messages")
	}

	return c.JSON(http.StatusOK, messages)
}

func (handler *MessengerHandler) getConferences(c echo.Context) error {
	userID, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
	if err != nil {
		handler.log.Error("Wrong ID type", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong ID type")
	}

	conferences, err := handler.service.GetConferencesByUser(userID)
	if err != nil {
		handler.log.Error("Failed to get conferences by user ID", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to get conferences")
	}

	return c.JSON(http.StatusOK, conferences)
}

func (handler *MessengerHandler) createConference(c echo.Context) error {
	var decoded createConferenceRequest
	err := json.NewDecoder(c.Request().Body).Decode(&decoded)
	if err != nil {
		handler.log.Error("Failed to create conference: parsing request body", zap.Error(err))
		return c.String(http.StatusBadRequest, "Need IDs of users and conference name")
	}

	err = handler.service.CreateConference(decoded.usersIDs, decoded.name, decoded.createdBy, decoded.createdAt)
	if err != nil {
		handler.log.Error("Failed to create conference", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to create conference")
	}
	return c.String(http.StatusCreated, "New conference created")
}
