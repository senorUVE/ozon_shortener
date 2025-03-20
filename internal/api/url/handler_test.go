package url_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"ozon_shortener/internal/api/url"
	"strings"
	"testing"

	valMocks "ozon_shortener/internal/middleware/validator/mocks"
	srvMocks "ozon_shortener/internal/services/url/mocks"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateURLHandlerOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvMocks.NewMockService(ctrl)
	mockValidator := valMocks.NewMockValidator(ctrl)

	originalURLs := []string{"https://example.com/page1"}
	expectedResult := map[string]string{
		"http://example.com/page1": "http://localhost:8080/abc12345_A",
	}

	mockValidator.EXPECT().ValidateURLs(originalURLs).Return(nil)
	mockService.EXPECT().CreateURL(gomock.Any(), originalURLs).Return(expectedResult, nil)

	handler := url.New(mockService, mockValidator)

	reqBody, err := json.Marshal(map[string][]string{"original_urls": originalURLs})
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/url/generate", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateURL(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Urls map[string]string `json:"urls"`
	}
	err = json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, resp.Urls)
}
func TestCreateURLHandlerInvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvMocks.NewMockService(ctrl)
	mockValidator := valMocks.NewMockValidator(ctrl)

	handler := url.New(mockService, mockValidator)

	req := httptest.NewRequest("POST", "/url/generate", bytes.NewReader([]byte("not a json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateURL(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid request body")
}

func TestCreateURLHandlerValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvMocks.NewMockService(ctrl)
	mockValidator := valMocks.NewMockValidator(ctrl)

	originalURLs := []string{"invalid_url"}
	mockValidator.EXPECT().ValidateURLs(originalURLs).Return(errors.New("validation error"))
	mockService.EXPECT().CreateURL(gomock.Any(), gomock.Any()).Times(0)
	mockService.EXPECT().PublicURL().Times(0)

	handler := url.New(mockService, mockValidator)

	reqBody, err := json.Marshal(map[string][]string{"original_urls": originalURLs})
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/url/generate", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateURL(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "validation error")
}

func TestGetOriginalHandlerOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvMocks.NewMockService(ctrl)
	mockValidator := valMocks.NewMockValidator(ctrl)

	shortURLs := []string{"http://localhost:8080/abc123"}
	expectedResult := map[string]string{
		"http://localhost:8080/abc12345_A": "https://example.com/page1",
	}

	mockValidator.EXPECT().ValidateURLs(shortURLs).Return(nil)
	mockService.EXPECT().GetOriginal(gomock.Any(), shortURLs).Return(expectedResult, nil)

	handler := url.New(mockService, mockValidator)

	reqBody, err := json.Marshal(map[string][]string{"short_urls": shortURLs})
	assert.NoError(t, err)
	req := httptest.NewRequest("GET", "/url/original", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.GetOriginal(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Urls map[string]string `json:"urls"`
	}
	err = json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, resp.Urls)
}

func TestRedirectToOriginalHandlerOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvMocks.NewMockService(ctrl)
	mockValidator := valMocks.NewMockValidator(ctrl)

	token := "abc123"
	mockService.EXPECT().PublicURL().Return("localhost:8080")
	mockValidator.EXPECT().ValidateKey(token).Return(nil)
	shortURL := "http://localhost:8080/abc123"
	expectedMapping := map[string]string{
		shortURL: "https://example.com/page1",
	}
	mockService.EXPECT().GetOriginal(gomock.Any(), []string{shortURL}).Return(expectedMapping, nil)

	handler := url.New(mockService, mockValidator)
	req := httptest.NewRequest("GET", "/"+token, nil)
	req = mux.SetURLVars(req, map[string]string{"token": token})
	rec := httptest.NewRecorder()

	handler.RedirectToOriginal(rec, req)

	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "https://example.com/page1", rec.Header().Get("Location"))
}

func TestCreateURLHandler_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvMocks.NewMockService(ctrl)
	mockValidator := valMocks.NewMockValidator(ctrl)

	tests := []struct {
		name            string
		request         any
		expectedCode    int
		expectedMessage string
	}{
		{
			name:            "Unknown field",
			request:         map[string]string{"unknown": "example"},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "Invalid request body",
		},
		{
			name:            "Empty body",
			request:         "",
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "Invalid request body",
		},
		{
			name:            "Too big URL",
			request:         map[string][]string{"original_urls": {strings.Repeat("a", 2050)}},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "URL must be maximum 2048 letters long",
		},
		{
			name:            "Empty URL",
			request:         map[string][]string{"original_urls": {""}},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "URL must be a valid URL string",
		},
		{
			name:            "Invalid URL #1",
			request:         map[string][]string{"original_urls": {"http"}},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "URL must be an absolute URL",
		},
		{
			name:            "Invalid URL #2",
			request:         map[string][]string{"original_urls": {"http://"}},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "URL must be an absolute URL",
		},
		{
			name:            "Invalid URL #3",
			request:         map[string][]string{"original_urls": {"httpss://example.com"}},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "URL must begin with http or https",
		},
		{
			name:            "Invalid URL #4",
			request:         map[string][]string{"original_urls": {"example.com"}},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "URL must be an absolute URL",
		},
		{
			name:            "Invalid URL #5",
			request:         map[string][]string{"original_urls": {"/example.com"}},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "URL must be an absolute URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rawBody []byte
			switch v := tt.request.(type) {
			case string:
				rawBody = []byte(v)
			default:
				var err error
				rawBody, err = json.Marshal(v)
				require.NoError(t, err)
			}

			switch tt.request.(type) {
			case map[string][]string:
				orig, _ := tt.request.(map[string][]string)["original_urls"]
				mockValidator.EXPECT().ValidateURLs(orig).Return(errors.New(tt.expectedMessage))
				mockService.EXPECT().CreateURL(gomock.Any(), gomock.Any()).Times(0)
				mockService.EXPECT().PublicURL().Times(0)
			}

			req := httptest.NewRequest("POST", "/url/generate", bytes.NewReader(rawBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := url.New(mockService, mockValidator)
			handler.CreateURL(rec, req)

			result := rec.Result()

			assert.Equal(t, tt.expectedCode, result.StatusCode)

			bodyBytes, err := io.ReadAll(result.Body)
			result.Body.Close()
			require.NoError(t, err)

			require.Contains(t, string(bodyBytes), tt.expectedMessage)
		})
	}
}

func TestGetOriginalHandlerBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvMocks.NewMockService(ctrl)
	mockValidator := valMocks.NewMockValidator(ctrl)

	tests := []struct {
		name         string
		request      any
		expectedCode int
		errorMessage string
	}{
		{
			name:         "Unknown field",
			request:      map[string]string{"unknown": "example"},
			expectedCode: http.StatusBadRequest,
			errorMessage: `json: unknown field "unknown"`,
		},
		{
			name:         "Invalid key",
			request:      map[string][]string{"short_urls": {"й"}},
			expectedCode: http.StatusBadRequest,
			errorMessage: "invalid letter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rawBody []byte
			switch v := tt.request.(type) {
			case string:
				rawBody = []byte(v)
			default:
				var err error
				rawBody, err = json.Marshal(v)
				require.NoError(t, err)
			}

			switch tt.request.(type) {
			case map[string][]string:
				shortURLs, _ := tt.request.(map[string][]string)["short_urls"]
				mockValidator.EXPECT().ValidateURLs(shortURLs).Return(errors.New(tt.errorMessage))
				mockService.EXPECT().GetOriginal(gomock.Any(), gomock.Any()).Times(0)
			}

			req := httptest.NewRequest("GET", "/url/original", bytes.NewReader(rawBody))
			//req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := url.New(mockService, mockValidator)
			handler.GetOriginal(rec, req)

			result := rec.Result()
			//require.Equal(t, "application/json", result.Header.Get("Content-Type"))
			require.Equal(t, tt.expectedCode, result.StatusCode)

			bodyBytes, err := io.ReadAll(result.Body)
			result.Body.Close()
			require.NoError(t, err)
			require.Contains(t, string(bodyBytes), tt.errorMessage)
		})
	}
}
