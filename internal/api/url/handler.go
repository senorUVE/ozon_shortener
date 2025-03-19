package url

import (
	"encoding/json"
	"fmt"
	"net/http"
	urladapter "ozon_shortener/internal/adapters/url"

	"ozon_shortener/internal/middleware/validator"

	"github.com/gorilla/mux"
)

type Handler struct {
	urlAdapter urladapter.UrlAdapters
	validator  *validator.Validator
}

func New(
	urlAdapter urladapter.UrlAdapters,
	validator *validator.Validator,
) *Handler {
	return &Handler{
		urlAdapter: urlAdapter,
		validator:  validator,
	}
}

// CreateURL создает короткую ссылку
// @Summary Создать короткую ссылку
// @Description Принимает оригинальные URL и возвращает их короткие версии
// @Produce json
// @Param original_urls body object{OriginalUrls=[]string} true "Массив оригинальных URL"
// @Success 200 {object} object{links=map[string]string} "Список оригинальных и сокращенных ссылок"
// @Failure 400 {object} object{error=string} "Неверный запрос"
// @Failure 500 {object} object{error=string} "Ошибка на сервере"
// @Router /url/generate [post]
func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req struct {
		OriginalUrls []string `json:"original_urls"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateURLs(req.OriginalUrls); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return

	}

	res, err := h.urlAdapter.GenerateUrls(ctx, req.OriginalUrls)
	if err != nil {
		http.Error(w, fmt.Sprintf("что-то пошло не так: %v", err), 500)
		return
	}

	resp := struct {
		Urls map[string]string `json:"urls"`
	}{
		Urls: res,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("что-то пошло не так: %v", err), 500)
		return
	}
	w.WriteHeader(200)
}

// GetOriginal обрабатывает запросы на получение оригинальных ссылок по сокращённым ссылкам
// @Summary Получить оригинальные ссылки
// @Description Принимает сокращённые ссылки и возвращает их оригинальные версии
// @Produce json
// @Param short_urls body object{ShortUrls=[]string} true "Массив сокращённых ссылок"
// @Success 200 {object} object{links=map[string]string} "Список сокращённых и оригинальных ссылок"
// @Failure 400 {object} object{error=string} "Неверный запрос"
// @Failure 500 {object} object{error=string} "Ошибка на сервере"
// @Router /url/original [get]
func (h *Handler) GetOriginal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		ShortUrls []string `json:"short_urls"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	res, err := h.urlAdapter.GetOriginal(ctx, req.ShortUrls)
	if err != nil {
		http.Error(w, fmt.Sprintf("что-то пошло не так: %v", err), 500)
		return
	}

	resp := struct {
		Urls map[string]string `json:"urls"`
	}{
		Urls: res,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("что-то пошло не так: %v", err), 500)
		return
	}

	w.WriteHeader(200)
}

func (h *Handler) RedirectToOriginal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := mux.Vars(r)["token"]

	if err := h.validator.ValidateKey(token); err != nil {
		http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusBadRequest)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.urlAdapter.PublicURL(), token)
	res, err := h.urlAdapter.GetOriginal(ctx, []string{shortURL})
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
