package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func NewRouter(h *SubscriptionHandler, log *logrus.Logger) chi.Router {
	r := chi.NewRouter()

	mw := NewMiddleware(log)

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(mw.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/subscriptions", h.Create)
		r.Get("/subscriptions", h.List)
		r.Get("/subscriptions/total", h.TotalCost)
		r.Get("/subscriptions/{id}", h.GetByID)
		r.Put("/subscriptions/{id}", h.Update)
		r.Delete("/subscriptions/{id}", h.Delete)
	})

	return r
}
