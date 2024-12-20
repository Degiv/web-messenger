package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/Degiv/web-messenger/internals/domain"
	"github.com/Degiv/web-messenger/internals/services/auth"
	mycookie "github.com/Degiv/web-messenger/pkg/cookie"
	"github.com/Degiv/web-messenger/pkg/pqerr"
)

const (
	IDCookieKey = "UserID"
)

type MessengerService interface {
	RegisterUser(username string, email string, password string) (int64, error)
	GetUserByID(userID int64) (*domain.User, error)
	VerifyConferenceMember(userID int64, conferenceID int64) (bool, error)
	GetConferenceMessages(conferenceID int64) ([]*domain.Message, error)
	GetConferencesByUser(userID int64) ([]*domain.Conference, error)
	CreateConference(usersIDs []int64, name string, createdBy int64, createdAt time.Time) error
	PostToConference(message *domain.Message) error
}

type AuthService interface {
	AuthorizeUser(username string, password string) (int64, error)
}

type Handler struct {
	messenger MessengerService
	auth      AuthService
	log       *zap.Logger
}

func NewMessengerHandler(messenger MessengerService, auth AuthService, log *zap.Logger) *Handler {
	return &Handler{
		messenger: messenger,
		auth:      auth,
		log:       log,
	}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) *echo.Echo {
	e.POST("login", h.login)
	e.POST("signUp", h.signUp)

	e.GET("messenger/conferences/:id", h.getMessages)
	e.POST("messenger/conferences/:id", h.postMessage)

	e.GET("messenger/conferences", h.getConferences)
	e.POST("messenger/conferences", h.createConference)

	e.GET("messenger/users/:id", h.getUser)
	return e
}

func (h *Handler) signUp(c echo.Context) error {
	var decoded signUpRequest
	err := json.NewDecoder(c.Request().Body).Decode(&decoded)
	if err != nil {
		h.log.Error("Bad request", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong user data")
	}

	id, err := h.messenger.RegisterUser(decoded.Username, decoded.Email, decoded.Password)

	if err != nil {
		if pqerr.IsUniqueViolatesError(err) {
			return c.String(http.StatusConflict, "User with this login or email already exists")
		}

		return c.String(http.StatusInternalServerError, "Server error")
	}

	const hoursNumber = 24
	cookie := mycookie.CreateCookie(
		IDCookieKey,
		strconv.FormatInt(id, 10),
		time.Now().Add(hoursNumber*time.Hour))

	c.SetCookie(cookie)
	return c.String(http.StatusOK, "Registered!")
}

func (h *Handler) login(c echo.Context) error {
	var decoded loginRequest

	err := json.NewDecoder(c.Request().Body).Decode(&decoded)
	if err != nil {
		h.log.Error("Bad request", zap.Error(err))
		return c.String(http.StatusBadRequest, "No username or password")
	}

	userID, err := h.auth.AuthorizeUser(decoded.Username, decoded.Password)
	switch {
	case errors.Is(err, auth.ErrNoSuchUser) || errors.Is(err, auth.ErrWrongPassword):
		return c.String(http.StatusUnauthorized, "Wrong username or password")
	case err != nil:
		h.log.Error("Failed to authorize user", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Server error")
	}

	const hoursNumber = 24
	cookie := mycookie.CreateCookie(
		IDCookieKey,
		strconv.FormatInt(userID, 10),
		time.Now().Add(hoursNumber*time.Hour))

	c.SetCookie(cookie)
	return c.String(http.StatusOK, "Authorized!")
}

func (h *Handler) getUser(c echo.Context) error {
	userIDStr := c.Param("id")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.log.Error("Wrong ID type", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong ID type")
	}

	user, err := h.messenger.GetUserByID(userID)
	if err != nil {
		h.log.Error("Failed to get user by ID", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to get user")
	}

	user.PasswordHash = ""
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) getMessages(c echo.Context) error {
	conferenceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Error("Wrong ID type", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong ID type")
	}

	userIDCookie, err := c.Cookie(IDCookieKey)
	if err != nil {
		h.log.Error("no cookie", zap.Error(err))
		return c.String(http.StatusUnauthorized, "Need to login")
	}

	userID, err := strconv.ParseInt(userIDCookie.Value, 10, 64)
	if err != nil {
		h.log.Error("Wrong ID type", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong ID type")
	}

	ok, err := h.messenger.VerifyConferenceMember(userID, conferenceID)
	if err != nil {
		h.log.Error("Failed to verify conference member", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Server error")
	}
	if !ok {
		return c.String(http.StatusUnauthorized, "You are not a member of conference")
	}

	messages, err := h.messenger.GetConferenceMessages(conferenceID)
	if err != nil {
		h.log.Error("Failed to get messages by conference ID", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to get messages")
	}

	return c.JSON(http.StatusOK, messages)
}

func (h *Handler) postMessage(c echo.Context) error {
	var message domain.Message
	err := json.NewDecoder(c.Request().Body).Decode(&message)
	if err != nil {
		h.log.Error("Failed to parse message", zap.Error(err))
		return c.String(http.StatusBadRequest, "Failed to post message: cannot parse message")
	}

	userIDCookie, err := c.Cookie(IDCookieKey)
	if err != nil {
		h.log.Error("no cookie")
		return c.String(http.StatusUnauthorized, "Need to login")
	}

	userID, err := strconv.ParseInt(userIDCookie.Value, 10, 64)
	if err != nil {
		h.log.Error("Failed to parse int from cookie", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong cookie type")
	}

	ok, err := h.messenger.VerifyConferenceMember(userID, message.ConferenceID)
	if err != nil {
		h.log.Error("Failed to verify conference member", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Server error")
	}
	if !ok {
		return c.String(http.StatusUnauthorized, "You are not a member of conference")
	}

	err = h.messenger.PostToConference(&message)
	if err != nil {
		h.log.Error("Failed to post message to conference", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to post message")
	}
	return c.String(http.StatusCreated, "Message posted")
}

func (h *Handler) getConferences(c echo.Context) error {
	cookie, err := c.Cookie(IDCookieKey)
	if err != nil {
		return err
	}

	userID, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		h.log.Error("Wrong ID type", zap.Error(err))
		return c.String(http.StatusBadRequest, "Wrong ID type")
	}

	conferences, err := h.messenger.GetConferencesByUser(userID)
	if err != nil {
		h.log.Error("Failed to get conferences by user ID", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to get conferences")
	}

	return c.JSON(http.StatusOK, conferences)
}

func (h *Handler) createConference(c echo.Context) error {
	var decoded createConferenceRequest
	err := json.NewDecoder(c.Request().Body).Decode(&decoded)
	if err != nil {
		h.log.Error("Failed to create conference: parsing request body", zap.Error(err))
		return c.String(http.StatusBadRequest, "Need IDs of users and conference name")
	}

	err = h.messenger.CreateConference(decoded.UsersIDs, decoded.Name, decoded.CreatedBy, decoded.CreatedAt)
	if err != nil {
		h.log.Error("Failed to create conference", zap.Error(err))
		return c.String(http.StatusInternalServerError, "Failed to create conference")
	}
	return c.String(http.StatusCreated, "New conference created")
}
