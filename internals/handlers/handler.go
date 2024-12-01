package handlers

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"messenger/internals/domain"
	"messenger/internals/services/auth"
	mycookie "messenger/pkg/cookie"
	"net/http"
	"strconv"
	"time"
)

const (
	IDCookieKey = "UserID"
)

type MessengerService interface {
	RegisterUser(username string, email string, password string) (bool, int64, error)
	GetUserByID(userID int64) (*domain.User, error)
	GetConferenceMessages(conferenceID int64) ([]*domain.Message, error)
	GetConferencesByUser(userID int64) ([]*domain.Conference, error)
	CreateConference(usersIDs []int64, name string, createdBy int64, createdAt time.Time) error
	PostToConference(message *domain.Message) error
}

type AuthService interface {
	AuthorizeUser(username string, password string) (int64, error)
}

type MessengerHandler struct {
	messenger MessengerService
	auth      AuthService
	log       *zap.Logger
}

func NewMessengerHandler(messenger MessengerService, auth AuthService, log *zap.Logger) *MessengerHandler {
	return &MessengerHandler{
		messenger: messenger,
		auth:      auth,
		log:       log,
	}
}

func (handler *MessengerHandler) RegisterRoutes(e *echo.Echo) *echo.Echo {
	e.GET("login", handler.login)
	e.POST("signIn", handler.signIn)

	e.GET("messenger/conferences/{id}", handler.getMessages)
	e.POST("messenger/conferences/{id}", handler.postMessage)

	e.GET("messenger/conferences", handler.getConferences)
	e.POST("messenger/conferences", handler.createConference)

	e.GET("messenger/users/{id}", handler.getUser)
	return e
}

func (handler *MessengerHandler) signIn(c echo.Context) error {
	var decoded signInRequest
	err := json.NewDecoder(c.Request().Body).Decode(&decoded)
	if err != nil {
		handler.log.Error("Bad request", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong user data")
	}

	ok, id, err := handler.messenger.RegisterUser(decoded.Username, decoded.Email, decoded.Password)
	if err != nil {
		handler.log.Error("Failed to sign in", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Server Error")
	}

	if !ok {
		return c.String(http.StatusConflict, "Username already exists")
	}

	cookie := mycookie.CreateCookie(
		IDCookieKey,
		strconv.FormatInt(id, 10),
		time.Now().Add(24*time.Hour))

	c.SetCookie(cookie)
	return c.String(http.StatusOK, "Registered!")
}

func (handler *MessengerHandler) login(c echo.Context) error {
	var decoded loginRequest

	err := json.NewDecoder(c.Request().Body).Decode(&decoded)
	if err != nil {
		handler.log.Error("Bad request", zap.Error(err))
		return c.String(http.StatusBadRequest, "No username or password")
	}

	userID, err := handler.auth.AuthorizeUser(decoded.Username, decoded.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrNoSuchUser) || errors.Is(err, auth.ErrWrongPassword):
			return c.String(http.StatusUnauthorized, "Wrong username or password")
		default:
			handler.log.Error("Failed to authorize user", zap.Error(err))
			return c.String(http.StatusInternalServerError, "Server error")
		}
	}

	cookie := mycookie.CreateCookie(
		IDCookieKey,
		strconv.FormatInt(userID, 10),
		time.Now().Add(24*time.Hour))

	c.SetCookie(cookie)
	return c.String(http.StatusOK, "Authorized!")
}

func (handler *MessengerHandler) getUser(c echo.Context) error {
	cookie, err := c.Cookie(IDCookieKey)
	if err != nil {
		return err
	}

	userID, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		handler.log.Error("Wrong ID type", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong ID type")
	}

	user, err := handler.messenger.GetUserByID(userID)
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

	messages, err := handler.messenger.GetConferenceMessages(conferenceID)
	if err != nil {
		handler.log.Error("Failed to get messages by conference ID", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to get messages")
	}

	return c.JSON(http.StatusOK, messages)
}

func (handler *MessengerHandler) postMessage(c echo.Context) error {
	var message *domain.Message
	err := json.NewDecoder(c.Request().Body).Decode(message)
	if err != nil {
		handler.log.Error("Failed to parse message", zap.Error(err))
		return c.String(http.StatusBadRequest, "Failed to post message: cannot parse message")
	}

	err = handler.messenger.PostToConference(message)
	if err != nil {
		handler.log.Error("Failed to post message to conference", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to post message")
	}
	return c.String(http.StatusCreated, "Message posted")
}

func (handler *MessengerHandler) getConferences(c echo.Context) error {
	cookie, err := c.Cookie(IDCookieKey)
	if err != nil {
		return err
	}

	userID, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		handler.log.Error("Wrong ID type", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong ID type")
	}
	if err != nil {
		handler.log.Error("Wrong ID type", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong ID type")
	}

	conferences, err := handler.messenger.GetConferencesByUser(userID)
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

	err = handler.messenger.CreateConference(decoded.UsersIDs, decoded.Name, decoded.CreatedBy, decoded.CreatedAt)
	if err != nil {
		handler.log.Error("Failed to create conference", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to create conference")
	}
	return c.String(http.StatusCreated, "New conference created")
}
