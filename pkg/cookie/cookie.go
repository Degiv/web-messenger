package cookie

import (
	"net/http"
	"time"
)

func CreateCookie(name string, value string, time time.Time) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = time
	return cookie
}
