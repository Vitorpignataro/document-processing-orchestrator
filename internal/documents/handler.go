package documents

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)


type Handler struct{
	service *Service
	repo 	*Repository
}

func NewHandler(
	service *Service,
	repo 	*Repository,
) *Handler {
	return &Handler{
		service: service,
		repo: repo,
	}
}

func (h *Handler) Create (
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var input CreateInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
        writeJSON(w, http.StatusBadRequest, map[string]string{
            "error": "request must contain one JSON object",
        })
        return
    }

	document, err := h.service.Create(r.Context(), input)

	if errors.Is(err, ErrInvalidInput) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
            "error": "invalid document metadata",
        })
        return
	}

	if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{
            "error": "failed to create document",
        })
        return
    }

	w.Header().Set("Location", "/documents/"+document.ID)
	writeJSON(w, http.StatusAccepted, document)
}

func (h *Handler) GetByID(
    w http.ResponseWriter,
    r *http.Request,
) {
    id := chi.URLParam(r, "id")

    document, err := h.repo.GetByID(r.Context(), id)

    if errors.Is(err, ErrNotFound) {
        writeJSON(w, http.StatusNotFound, map[string]string{
            "error": "document not found",
        })
        return
    }

    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{
            "error": "failed to retrieve document",
        })
        return
    }

    writeJSON(w, http.StatusOK, document)
}


func writeJSON(
    w http.ResponseWriter,
    status int,
    value any,
) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(value)
}
