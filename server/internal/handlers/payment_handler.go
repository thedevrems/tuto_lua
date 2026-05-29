package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/thedevrems/tuto_lua/server/internal/auth"
	"github.com/thedevrems/tuto_lua/server/internal/httpx"
	"github.com/thedevrems/tuto_lua/server/internal/payment"
	"github.com/thedevrems/tuto_lua/server/internal/store"
)

// maxWebhookBody caps the Stripe webhook payload size.
const maxWebhookBody = 1 << 16

// PaymentHandler exposes Stripe Checkout creation and the Stripe webhook.
type PaymentHandler struct {
	svc *payment.Service
}

// NewPaymentHandler builds the handler around the payment service.
func NewPaymentHandler(svc *payment.Service) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

type checkoutRequest struct {
	CourseID string `json:"courseId"`
}

// Checkout creates a Stripe Checkout session and returns its redirect URL.
func (h *PaymentHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFrom(r.Context())
	var req checkoutRequest
	if !bind(w, r, &req) {
		return
	}
	url, err := h.svc.CreateCheckout(user.ID, req.CourseID)
	if err != nil {
		writeCheckoutError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"url": url})
}

// Webhook receives Stripe events and unlocks courses on successful payment.
// The raw body is required for signature verification.
func (h *PaymentHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxWebhookBody))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "payload illisible")
		return
	}
	if err := h.svc.HandleWebhook(body, r.Header.Get("Stripe-Signature")); err != nil {
		httpx.Error(w, http.StatusBadRequest, "webhook invalide")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// writeCheckoutError maps payment-service errors to HTTP statuses.
func writeCheckoutError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, payment.ErrNotConfigured):
		httpx.Error(w, http.StatusServiceUnavailable, "le paiement n'est pas configuré sur ce serveur")
	case errors.Is(err, payment.ErrFreeCourse):
		httpx.Error(w, http.StatusBadRequest, "ce cours est gratuit")
	case errors.Is(err, payment.ErrAlreadyEnrolled):
		httpx.Error(w, http.StatusConflict, "vous avez déjà accès à ce cours")
	case errors.Is(err, store.ErrNotFound):
		httpx.Error(w, http.StatusNotFound, "cours introuvable")
	default:
		httpx.Error(w, http.StatusInternalServerError, "création du paiement impossible")
	}
}
