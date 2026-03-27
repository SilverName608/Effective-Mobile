package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/SilverName608/Effective-Mobile/internal/model"
	"github.com/SilverName608/Effective-Mobile/internal/service"
)

type SubscriptionHandler struct {
	svc service.SubscriptionServiceI
	log *logrus.Logger
}

func NewSubscriptionHandler(svc service.SubscriptionServiceI, log *logrus.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc, log: log}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateSubscriptionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.log.WithError(err).Warn("invalid create request body")
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid start_date, expected MM-YYYY")
		return
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
	}

	if req.EndDate != nil {
		endDate, err := parseMonthYear(*req.EndDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid end_date, expected MM-YYYY")
			return
		}
		sub.EndDate = &endDate
	}

	if err := h.svc.Create(r.Context(), sub); err != nil {
		h.log.WithError(err).Error("create subscription failed")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, sub)
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	sub, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.log.WithError(err).Error("get subscription failed")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if sub == nil {
		writeError(w, http.StatusNotFound, "subscription not found")
		return
	}

	writeJSON(w, http.StatusOK, sub)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	var userID *uuid.UUID
	var serviceName *string

	v := r.URL.Query().Get("user_id")
	if v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		userID = &id
	}

	v = r.URL.Query().Get("service_name")
	if v != "" {
		serviceName = &v
	}

	subs, err := h.svc.List(r.Context(), userID, serviceName)
	if err != nil {
		h.log.WithError(err).Error("list subscriptions failed")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if subs == nil {
		subs = []*model.Subscription{}
	}

	writeJSON(w, http.StatusOK, subs)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateSubscriptionRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sub, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if sub == nil {
		writeError(w, http.StatusNotFound, "subscription not found")
		return
	}

	if req.ServiceName != nil {
		sub.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		sub.Price = *req.Price
	}
	if req.UserID != nil {
		uid, err := uuid.Parse(*req.UserID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		sub.UserID = uid
	}
	if req.StartDate != nil {
		t, err := parseMonthYear(*req.StartDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid start_date, expected MM-YYYY")
			return
		}
		sub.StartDate = t
	}
	if req.EndDate != nil {
		t, err := parseMonthYear(*req.EndDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid end_date, expected MM-YYYY")
			return
		}
		sub.EndDate = &t
	}

	err = h.svc.Update(r.Context(), sub)
	if err != nil {
		h.log.WithError(err).Error("update subscription failed")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, sub)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.svc.Delete(r.Context(), id)
	if err != nil {
		h.log.WithError(err).Error("delete subscription failed")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) TotalCost(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		writeError(w, http.StatusBadRequest, "from and to are required")
		return
	}

	from, err := parseMonthYear(fromStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid from, expected MM-YYYY")
		return
	}
	to, err := parseMonthYear(toStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid to, expected MM-YYYY")
		return
	}

	var userID *uuid.UUID
	var serviceName *string

	v := r.URL.Query().Get("user_id")
	if v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		userID = &id
	}

	v = r.URL.Query().Get("service_name")
	if v != "" {
		serviceName = &v
	}

	total, err := h.svc.TotalCost(r.Context(), userID, serviceName, from, to)
	if err != nil {
		h.log.WithError(err).Error("total cost failed")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, TotalCostResponse{Total: total})
}

func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
