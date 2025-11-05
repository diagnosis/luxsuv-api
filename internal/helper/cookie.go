package helper

import (
	"net/http"
	"time"
)

const RefreshCookieName = "refresh_token"

func SetRefreshCookie(w http.ResponseWriter, value string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshCookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   getEnvDef("APP_ENV", "development") == "production",
		MaxAge:   int(ttl.Seconds()),
	})
}

func ClearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   getEnvDef("APP_ENV", "development") == "production",
		MaxAge:   -1,
	})
}
