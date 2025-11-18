package cookies

import (
	"net/http"

	"github.com/prabalesh/loco/backend/pkg/config"
)

type CookieManager struct {
	cfg *config.Config
}

func NewCookieManager(cfg *config.Config) *CookieManager {
	return &CookieManager{cfg: cfg}
}

func (cm *CookieManager) SetSecure(w http.ResponseWriter, name, value string, maxAge int) {
	sameSite := http.SameSiteLaxMode
	switch cm.cfg.Cookie.SameSite {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	case "lax":
		sameSite = http.SameSiteLaxMode
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   cm.cfg.Cookie.Secure,
		SameSite: sameSite,
		Domain:   cm.cfg.Cookie.Domain,
	}
	http.SetCookie(w, cookie)
}

func (cm *CookieManager) Clear(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cm.cfg.Cookie.Secure,
		Domain:   cm.cfg.Cookie.Domain,
	}
	http.SetCookie(w, cookie)
}
