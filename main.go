package secret_auth

import (
	"context"
	"fmt"
	"net/http"
)

type Config struct {
	CookieName string `json:"cookieName,omitempty"`
	SecretKey  string `json:"secretKey,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		CookieName: "secret",
		SecretKey:  "",
	}
}

type SecretAuthPlugin struct {
	next       http.Handler
	cookieName string
	secretKey  string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.SecretKey) == 0 {
		return nil, fmt.Errorf("secret key cannot be empty")
	}

	return &SecretAuthPlugin{
		next:       next,
		cookieName: config.CookieName,
		secretKey:  config.SecretKey,
	}, nil
}

func (a *SecretAuthPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(a.cookieName)

	if err != nil || cookie.Value != a.secretKey {
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}

	a.next.ServeHTTP(rw, req)
}
