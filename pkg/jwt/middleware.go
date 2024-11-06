package jwt

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// TODO: replace to config
const secret = "secret"

func Authorization(next echo.HandlerFunc) echo.HandlerFunc {
	return echojwt.JWT([]byte(secret))(next)
}
