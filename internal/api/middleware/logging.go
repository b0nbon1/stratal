package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	reset  = "\033[0m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	cyan   = "\033[36m"
	bold   = "\033[1m"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	n, err := sr.ResponseWriter.Write(b)
	sr.size += n
	return n, err
}

const slowThreshold = 500 * time.Millisecond // Customize as needed

type ctxKeyRequestID struct{}

func Logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.NewString()
			}

			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			}
			method := r.Method
			path := r.URL.Path
			query := ""
			if r.URL.RawQuery != "" {
				query = "?" + r.URL.RawQuery
			}
			timestamp := time.Now().Format("15:04:05")

			fmt.Printf("[%s] 🚀 %s%s %s from %s | ID: %s\n",
				yellow+timestamp+reset, bold+cyan, method, green+path+query+reset, ip, requestID)

			// Track request
			rec := &statusRecorder{ResponseWriter: w, status: 200}
			r = r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID{}, requestID))
			next.ServeHTTP(rec, r)

			duration := time.Since(start)
			logStatus := colorStatus(rec.status)
			logSize := fmt.Sprintf("%s%dB%s", cyan, rec.size, reset)

			msg := fmt.Sprintf(" - %s%s (%s) from %s | ID: %s | Size: %s | Time: %v%s\n",
				logStatus, path+query, method, ip, requestID, logSize, duration, reset,
			)

			if duration > slowThreshold {
				fmt.Printf("                     %s⚠️  SLOW%s%s", red, reset, msg)
			} else {
				fmt.Printf("                     ✅ %sCompleted%s%s", bold, reset, msg)
			}
		})
	}
}

func colorStatus(status int) string {
	code := fmt.Sprintf("%d", status)
	switch {
	case status >= 200 && status < 300:
		return green + code + reset + " "
	case status >= 400 && status < 500:
		return yellow + code + reset + " "
	case status >= 500:
		return red + code + reset + " "
	default:
		return code + " "
	}
}
