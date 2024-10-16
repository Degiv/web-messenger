package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	domain2 "messenger/internals/domain"
	"net/http"
	"strconv"
	"time"
)

type MessengerService interface {
	GetUserByID(userID int64) (*domain2.User, error)
	GetConferenceMessages(conferenceID int64) ([]*domain2.Message, error)
	GetConferencesByUser(userID int64) ([]*domain2.Conference, error)
	CreateConference(usersIDs []int64, name string, createdBy int64, createdAt time.Time) error
	PostToConference(message *domain2.Message) error
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
	e.POST("messenger/conferences/{id}", handler.postMessage)

	e.GET("messenger/conferences", handler.getConferences)
	e.POST("messenger/conferences", handler.createConference)

	e.GET("messenger/users/{id}", handler.getUser)
	//e.POST("messenger/users", handler.createUser)
	return e
}

func (handler *MessengerHandler) getUser(c echo.Context) error {
	userID, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
	if err != nil {
		handler.log.Error("Wrong ID type", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong ID type")
	}

	user, err := handler.service.GetUserByID(userID)
	if err != nil {
		handler.log.Error("Failed to get user by ID", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to get user")
	}

	return c.JSON(http.StatusOK, user)
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

func (handler *MessengerHandler) postMessage(c echo.Context) error {
	var message *domain2.Message
	err := json.NewDecoder(c.Request().Body).Decode(message)
	if err != nil {
		handler.log.Error("Failed to parse message", zap.Error(err))
		return c.String(http.StatusBadRequest, "Failed to post message: cannot parse message")
	}

	err = handler.service.PostToConference(message)
	if err != nil {
		handler.log.Error("Failed to post message to conference", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to post message")
	}
	return c.String(http.StatusCreated, "Message posted")
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

	err = handler.service.CreateConference(decoded.UsersIDs, decoded.Name, decoded.CreatedBy, decoded.CreatedAt)
	if err != nil {
		handler.log.Error("Failed to create conference", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to create conference")
	}
	return c.String(http.StatusCreated, "New conference created")
}
