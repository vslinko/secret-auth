package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is starting...")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Extract the username and password from the Authorization header
	user, pass, ok := parseAuth(r)
	if !ok {
		// If parsing fails or credentials are not provided, request authentication
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// In a real application, you should validate the username and password
	// Here we are just checking against hardcoded credentials
	if user != os.Getenv("S_USERNAME") || pass != os.Getenv("S_PASSWORD") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// If authentication is successful, set a cookie
	http.SetCookie(w, &http.Cookie{
		Name:   os.Getenv("S_COOKIE_NAME"),
		Value:  os.Getenv("S_COOKIE_VALUE"),
		Path:   "/",
		Domain: os.Getenv("S_COOKIE_DOMAIN"),
		MaxAge: 60 * 60 * 24 * 365,
	})

	// Respond to the client; check for the return_url query param
	// Use the return_url for redirection if present, otherwise use the environment variable
	redirectURL := os.Getenv("S_REDIRECT_URL") // Default redirect URL from env var
	if returnURL := r.URL.Query().Get("return_url"); returnURL != "" {
		redirectURL = returnURL // Override with the return_url query param if present
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func parseAuth(r *http.Request) (username, password string, ok bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", "", false
	}

	// Verify that authentication is Basic
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return "", "", false
	}

	// Decode the base64 encoded credentials
	cred, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return "", "", false
	}

	// Split username and password
	parts := strings.SplitN(string(cred), ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	return parts[0], parts[1], true
}
