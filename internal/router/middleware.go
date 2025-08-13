package router

import (
	"net/http"
	"net/url"
	"os"
	"regexp"
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
		hostname := originURL.Hostname()
		for _, domain := range strings.Split(allowedDomains, ",") {
			trimmedDomain := strings.TrimSpace(domain)
			if strings.Contains(trimmedDomain, "*") {
				pattern := strings.ReplaceAll(trimmedDomain, ".", "\\.")
				pattern = strings.ReplaceAll(pattern, "*", "[^.]+")
				pattern = "^" + pattern + "$"
				if re, err := regexp.Compile(pattern); err == nil {
					if re.MatchString(hostname) {
						isAllowed = true
						break
					}
				}
			} else if trimmedDomain == hostname {
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
