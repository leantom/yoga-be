package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
	"github.com/quangho/yoga-be/internal/adapter/firestoredb"
)

type resourceConfig struct {
	Path        string
	Collection  string
	Filters     map[string]string
	Orders      []string
	DefaultSort []string
	CreateNow   []string
	UpdateNow   []string
}

type resourceHandler struct {
	repo firestoredb.Repository
	cfg  resourceConfig
}

func mountResource(r chi.Router, registry firestoredb.Registry, cfg resourceConfig) {
	h := resourceHandler{repo: registry.Repository(cfg.Collection), cfg: cfg}
	r.Route(cfg.Path, func(r chi.Router) {
		r.Get("", h.list)
		r.Get("/", h.list)
		r.Post("", h.create)
		r.Post("/", h.create)
		r.Get("/{id}", h.get)
		r.Patch("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
}

func (h resourceHandler) create(w http.ResponseWriter, r *http.Request) {
	body, ok := decodeObject(w, r)
	if !ok {
		return
	}
	setMissingTimes(body, h.cfg.CreateNow)
	item, err := h.repo.Create(r.Context(), body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h resourceHandler) get(w http.ResponseWriter, r *http.Request) {
	item, err := h.repo.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h resourceHandler) update(w http.ResponseWriter, r *http.Request) {
	body, ok := decodeObject(w, r)
	if !ok {
		return
	}
	setTimes(body, h.cfg.UpdateNow)
	item, err := h.repo.Update(r.Context(), chi.URLParam(r, "id"), body)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h resourceHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		handleRepoError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h resourceHandler) list(w http.ResponseWriter, r *http.Request) {
	limit, err := firestoredb.ParseLimit(r.URL.Query().Get("limit"), 50, 200)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	opts := firestoredb.ListOptions{Limit: limit}
	for field, op := range h.cfg.Filters {
		raw := r.URL.Query().Get(field)
		if raw == "" {
			continue
		}
		value, err := coerceQueryValue(raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		opts.Filters = append(opts.Filters, firestoredb.Filter{Field: field, Op: op, Value: value})
	}

	sortParts := h.cfg.DefaultSort
	if rawSort := r.URL.Query().Get("sort"); rawSort != "" {
		sortParts = strings.Split(rawSort, ",")
	}
	for _, part := range sortParts {
		field, direction, err := parseSort(part)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if !slices.Contains(h.cfg.Orders, field) {
			writeError(w, http.StatusBadRequest, errors.New("unsupported sort field: "+field))
			return
		}
		opts.OrderBy = append(opts.OrderBy, firestoredb.Order{Field: field, Direction: direction})
	}

	items, err := h.repo.List(r.Context(), opts)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "limit": limit})
}

func decodeObject(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	defer r.Body.Close()
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return nil, false
	}
	if body == nil {
		writeError(w, http.StatusBadRequest, errors.New("request body must be a JSON object"))
		return nil, false
	}
	return body, true
}

func parseSort(part string) (string, firestore.Direction, error) {
	pieces := strings.Split(strings.TrimSpace(part), ":")
	field := pieces[0]
	direction := firestore.Asc
	if len(pieces) > 1 {
		switch strings.ToLower(pieces[1]) {
		case "asc":
			direction = firestore.Asc
		case "desc":
			direction = firestore.Desc
		default:
			return "", firestore.Asc, errors.New("sort direction must be asc or desc")
		}
	}
	return field, direction, nil
}

func coerceQueryValue(raw string) (any, error) {
	if value, err := firestoredb.ParseBool(raw); err == nil {
		return value, nil
	}
	if value, err := strconv.Atoi(raw); err == nil {
		return value, nil
	}
	if value, err := strconv.ParseFloat(raw, 64); err == nil {
		return value, nil
	}
	return raw, nil
}

func setMissingTimes(body map[string]any, fields []string) {
	now := time.Now().UTC()
	for _, field := range fields {
		if _, exists := body[field]; !exists {
			body[field] = now
		}
	}
}

func setTimes(body map[string]any, fields []string) {
	now := time.Now().UTC()
	for _, field := range fields {
		body[field] = now
	}
}
