package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/JellyGoFF/FF-Hexas/core/breaker"
	"github.com/JellyGoFF/FF-Hexas/core/logc"
	"github.com/JellyGoFF/FF-Hexas/core/stat"
	"github.com/JellyGoFF/FF-Hexas/rest/httpx"
	"github.com/JellyGoFF/FF-Hexas/rest/internal/response"
)

const breakerSeparator = "://"

// BreakerHandler returns a break circuit middleware.
func BreakerHandler(method, path string, metrics *stat.Metrics) func(http.Handler) http.Handler {
	brk := breaker.NewBreaker(breaker.WithName(strings.Join([]string{method, path}, breakerSeparator)))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			promise, err := brk.Allow()
			if err != nil {
				metrics.AddDrop()
				logc.Errorf(r.Context(), "[http] dropped, %s - %s - %s",
					r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent())
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			cw := response.NewWithCodeResponseWriter(w)
			defer func() {
				if cw.Code < http.StatusInternalServerError {
					promise.Accept()
				} else {
					promise.Reject(fmt.Sprintf("%d %s", cw.Code, http.StatusText(cw.Code)))
				}
			}()
			next.ServeHTTP(cw, r)
		})
	}
}
