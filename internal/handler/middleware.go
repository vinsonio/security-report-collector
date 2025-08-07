package handler

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

// CORSMiddleware validates the request's Origin or Referer header against a whitelist of allowed domains.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedDomains := os.Getenv("ALLOWED_DOMAINS")
		if allowedDomains == "" {
			next.ServeHTTP(w, r)
			return
		}

		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = r.Header.Get("Referer")
		}

		if origin == "" {
			http.Error(w, "missing origin or referer header", http.StatusBadRequest)
			return
		}

		originURL, err := url.Parse(origin)
		if err != nil {
			http.Error(w, "invalid origin or referer header", http.StatusBadRequest)
			return
		}

		isAllowed := false
		for _, domain := range strings.Split(allowedDomains, ",") {
			if strings.TrimSpace(domain) == originURL.Hostname() {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			http.Error(w, "origin not allowed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}