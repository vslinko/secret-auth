package secret_auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Config struct {
	CookieName     string `json:"cookieName,omitempty"`
	HeaderName     string `json:"headerName,omitempty"`
	ParamName      string `json:"paramName,omitempty"`
	SecretKey      string `json:"secretKey,omitempty"`
	AuthUrl        string `json:"authUrl,omitempty"`
	ReturnUrlParam string `json:"returnUrlParam,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		CookieName:     "secret",
		HeaderName:     "",
		ParamName:      "",
		SecretKey:      "",
		AuthUrl:        "",
		ReturnUrlParam: "return_url",
	}
}

type SecretAuthPlugin struct {
	next           http.Handler
	cookieName     string
	headerName     string
	paramName      string
	secretKey      string
	authUrl        string
	returnUrlParam string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.SecretKey) == 0 {
		return nil, fmt.Errorf("secret key cannot be empty")
	}

	return &SecretAuthPlugin{
		next:           next,
		cookieName:     config.CookieName,
		headerName:     config.HeaderName,
		paramName:      config.ParamName,
		secretKey:      config.SecretKey,
		authUrl:        config.AuthUrl,
		returnUrlParam: config.ReturnUrlParam,
	}, nil
}

func (a *SecretAuthPlugin) checkAuth(req *http.Request) bool {
	cookie, err := req.Cookie(a.cookieName)

	if err == nil && cookie.Value == a.secretKey {
		return true
	}

	if a.headerName != "" {
		if header := req.Header.Get(a.headerName); header == a.secretKey {
			return true
		}
	}

	if a.paramName != "" {
		if param := req.URL.Query().Get(a.paramName); param == a.secretKey {
			return true
		}
	}

	return false
}

func (a *SecretAuthPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !a.checkAuth(req) {
		if a.authUrl != "" {
			// Obtain the complete request URL
			requestURL := getFullURL(req)

			// Construct the URL for redirection
			redirectURL := fmt.Sprintf("%s?%s=%s", a.authUrl, a.returnUrlParam, url.QueryEscape(requestURL))

			// Perform the redirection
			http.Redirect(rw, req, redirectURL, http.StatusTemporaryRedirect)
		} else {
			http.Error(rw, "Forbidden", http.StatusForbidden)
		}

		return
	}

	a.next.ServeHTTP(rw, req)
}

func getFullURL(req *http.Request) string {
	scheme := "http" // Default to HTTP
	if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
		// This checks if the underlying connection is TLS (indicating HTTPS) or
		// if the request was originally received as HTTPS before being forwarded (common in proxies/load balancers)
		scheme = "https"
	}

	fullURL := fmt.Sprintf("%s://%s%s", scheme, req.Host, req.URL.Path)

	// If there are query parameters, append them as well
	if rawQuery := req.URL.RawQuery; rawQuery != "" {
		fullURL += "?" + rawQuery
	}

	return fullURL
}
