package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Middleware struct {
	log *logrus.Logger
}

func NewMiddleware(log *logrus.Logger) *Middleware {
	return &Middleware{log: log}
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("incoming request")

		next.ServeHTTP(w, r)
	})
}
