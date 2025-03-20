package url

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ozon_shortener/internal/middleware/validator"
	"ozon_shortener/internal/services/url"

	"github.com/gorilla/mux"
)

type Handler struct {
	service   url.Service
	validator validator.Validator
}

func New(service url.Service, validator validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

// CreateURL создает короткие ссылки
// @Summary Создать короткую ссылку
// @Description Принимает оригинальные URL и возвращает их короткие версии
// @Produce json
// @Param original_urls body object{original_urls=[]string} true "Массив оригинальных URL"
// @Success 200 {object} object{urls=map[string]string} "Список оригинальных и сокращенных ссылок"
// @Failure 400 {object} object{error=string} "Неверный запрос"
// @Failure 500 {object} object{error=string} "Ошибка на сервере"
// @Router /url/generate [post]
func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req struct {
		OriginalUrls []string `json:"original_urls"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateURLs(req.OriginalUrls); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	res, err := h.service.CreateURL(ctx, req.OriginalUrls)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal error: %v", err), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Urls map[string]string `json:"urls"`
	}{
		Urls: res,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Internal error: %v", err), http.StatusInternalServerError)
		return
	}
}

// GetOriginal обрабатывает запросы на получение оригинальных ссылок по сокращённым ссылкам
// @Summary Получить оригинальные ссылки
// @Description Принимает абсолютные короткие ссылки и возвращает их оригинальные версии
// @Produce json
// @Param short_urls body object{short_urls=[]string} true "Массив коротких ссылок (абсолютный URL)"
// @Success 200 {object} object{urls=map[string]string} "Список сокращённых и оригинальных ссылок"
// @Failure 400 {object} object{error=string} "Неверный запрос"
// @Failure 500 {object} object{error=string} "Ошибка на сервере"
// @Router /url/original [get]
func (h *Handler) GetOriginal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req struct {
		ShortUrls []string `json:"short_urls"`
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateURLs(req.ShortUrls); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	res, err := h.service.GetOriginal(ctx, req.ShortUrls)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal error: %v", err), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Urls map[string]string `json:"urls"`
	}{
		Urls: res,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Internal error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := mux.Vars(r)["token"]

	if err := h.validator.ValidateKey(token); err != nil {
		http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusBadRequest)
		return
	}
	shortURL := fmt.Sprintf("http://%s/%s", h.service.PublicURL(), token)
	res, err := h.service.GetOriginal(ctx, []string{shortURL})
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	originalURL, ok := res[shortURL]
	if !ok || originalURL == "" {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
