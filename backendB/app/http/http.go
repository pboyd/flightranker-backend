package http

import (
	"net/http"

	"github.com/pboyd/flights/backend/backendb/app/graphql"
)

type Handler struct {
	Processor       *graphql.Processor
	CORSAllowOrigin string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.CORSAllowOrigin != "" {
		w.Header().Set("Access-Control-Allow-Origin", h.CORSAllowOrigin)
	}

	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	results, err := h.Processor.Do(r.Context(), query)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Write([]byte(results))
	w.Write([]byte("\n"))
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	if qe, ok := err.(graphql.QueryError); ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(qe.Error()))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`[{"message":"An internal error occurred"}]`))
	}
	w.Write([]byte("\n"))
}
