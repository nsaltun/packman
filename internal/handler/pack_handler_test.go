package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nsaltun/packman/internal/apperror"
	"github.com/nsaltun/packman/internal/middleware"
	"github.com/nsaltun/packman/internal/mocks"
	"github.com/nsaltun/packman/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	// Set Gin to test mode to avoid unnecessary output
	gin.SetMode(gin.TestMode)
}

// setupTestRouter creates a Gin engine with necessary middleware for testing
func setupTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.ErrorHandler())
	return router
}

func TestPackHTTPHandler_CalculatePacks(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*mocks.MockPackService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful calculation",
			requestBody: model.PackCalculationRequest{
				Quantity: 250,
			},
			mockSetup: func(m *mocks.MockPackService) {
				m.On("CalculatePacks", mock.Anything, 250).
					Return(&model.PackCalculationResponse{
						Quantity: 250,
						Packs:    map[int]int{250: 1},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, float64(250), data["quantity"])

				packs, ok := data["packs"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, float64(1), packs["250"])
			},
		},
		{
			name: "validation error - zero quantity",
			requestBody: model.PackCalculationRequest{
				Quantity: 0,
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				errorData, ok := response["error"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, string(apperror.ErrCodeValidation), errorData["code"])
			},
		},
		{
			name: "validation error - negative quantity",
			requestBody: model.PackCalculationRequest{
				Quantity: -10,
			},
			mockSetup: func(m *mocks.MockPackService) {
				// no mock calls expected
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				errorData, ok := response["error"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, string(apperror.ErrCodeValidation), errorData["code"])
			},
		},
		{
			name: "service error - not found",
			requestBody: model.PackCalculationRequest{
				Quantity: 100,
			},
			mockSetup: func(m *mocks.MockPackService) {
				m.On("CalculatePacks", mock.Anything, 100).
					Return(nil, apperror.NotFoundError("Pack configuration not found", nil))
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				errorData, ok := response["error"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, string(apperror.ErrCodeNotFound), errorData["code"])
			},
		},
		{
			name: "service error - internal error",
			requestBody: model.PackCalculationRequest{
				Quantity: 100,
			},
			mockSetup: func(m *mocks.MockPackService) {
				m.On("CalculatePacks", mock.Anything, 100).
					Return(nil, apperror.InternalError("Database error", errors.New("connection failed")))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				errorData, ok := response["error"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, string(apperror.ErrCodeInternal), errorData["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			mockService := new(mocks.MockPackService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}
			handler := NewPackHTTPHandler(mockService)

			// create request
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// create router with error handler middleware and execute
			w := httptest.NewRecorder()
			router := setupTestRouter()
			router.POST("/api/v1/calculate", handler.CalculatePacks)
			router.ServeHTTP(w, req)

			// assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			mockService.AssertExpectations(t)
		})
	}
}

func TestPackHTTPHandler_GetPackSizes(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockPackService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful retrieval",
			mockSetup: func(m *mocks.MockPackService) {
				m.On("GetPackSizes", mock.Anything).
					Return(&model.GetPackSizesResponse{
						PackSizes: []int{250, 500, 1000, 2000, 5000},
						Version:   1,
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedBy: "system",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok)

				packSizes, ok := data["pack_sizes"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, packSizes, 5)
				assert.Equal(t, float64(1), data["version"])
			},
		},
		{
			name: "service error - not found",
			mockSetup: func(m *mocks.MockPackService) {
				m.On("GetPackSizes", mock.Anything).
					Return(nil, apperror.NotFoundError("Pack configuration not found", nil))
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				errorData, ok := response["error"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, string(apperror.ErrCodeNotFound), errorData["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			mockService := new(mocks.MockPackService)
			tt.mockSetup(mockService)
			handler := NewPackHTTPHandler(mockService)

			// create request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/pack-sizes", nil)

			// create router with error handler middleware and execute
			w := httptest.NewRecorder()
			router := setupTestRouter()
			router.GET("/api/v1/pack-sizes", handler.GetPackSizes)
			router.ServeHTTP(w, req)

			// assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			mockService.AssertExpectations(t)
		})
	}
}

func TestPackHTTPHandler_UpdatePackSizes(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*mocks.MockPackService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful update",
			requestBody: model.UpdatePackSizesRequest{
				PackSizes: []int{250, 500, 1000},
				UpdatedBy: "admin",
			},
			mockSetup: func(m *mocks.MockPackService) {
				m.On("UpdatePackSizes", mock.Anything, []int{250, 500, 1000}, "admin").
					Return(&model.UpdatePackSizesResponse{
						PackSizes: []int{250, 500, 1000},
						Version:   2,
						UpdatedAt: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
						UpdatedBy: "admin",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok)

				packSizes, ok := data["pack_sizes"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, packSizes, 3)
				assert.Equal(t, float64(2), data["version"])
				assert.Equal(t, "admin", data["updated_by"])
			},
		},
		{
			name: "successful update with duplicates removed",
			requestBody: model.UpdatePackSizesRequest{
				PackSizes: []int{250, 500, 250, 1000, 500},
				UpdatedBy: "admin",
			},
			mockSetup: func(m *mocks.MockPackService) {
				// After deduplication, should be [250, 500, 1000]
				m.On("UpdatePackSizes", mock.Anything, []int{250, 500, 1000}, "admin").
					Return(&model.UpdatePackSizesResponse{
						PackSizes: []int{250, 500, 1000},
						Version:   2,
						UpdatedAt: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
						UpdatedBy: "admin",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok)

				packSizes, ok := data["pack_sizes"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, packSizes, 3)
				assert.Equal(t, float64(2), data["version"])
				assert.Equal(t, "admin", data["updated_by"])
			},
		},
		{
			name:           "invalid JSON request",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				errorData, ok := response["error"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, string(apperror.ErrCodeBadRequest), errorData["code"])
			},
		},
		{
			name: "validation error - empty pack sizes",
			requestBody: model.UpdatePackSizesRequest{
				PackSizes: []int{},
				UpdatedBy: "admin",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				errorData, ok := response["error"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, string(apperror.ErrCodeValidation), errorData["code"])
			},
		},
		{
			name: "validation error - zero pack size",
			requestBody: model.UpdatePackSizesRequest{
				PackSizes: []int{250, 0, 1000},
				UpdatedBy: "admin",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				errorData, ok := response["error"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, string(apperror.ErrCodeValidation), errorData["code"])
			},
		},
		{
			name: "service error - internal error",
			requestBody: model.UpdatePackSizesRequest{
				PackSizes: []int{250, 500},
				UpdatedBy: "admin",
			},
			mockSetup: func(m *mocks.MockPackService) {
				m.On("UpdatePackSizes", mock.Anything, []int{250, 500}, "admin").
					Return(nil, apperror.InternalError("Database error", errors.New("connection failed")))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				errorData, ok := response["error"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, string(apperror.ErrCodeInternal), errorData["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			mockService := new(mocks.MockPackService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}
			handler := NewPackHTTPHandler(mockService)

			// create request
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/pack-sizes", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// create router with error handler middleware and execute
			w := httptest.NewRecorder()
			router := setupTestRouter()
			router.PUT("/api/v1/pack-sizes", handler.UpdatePackSizes)
			router.ServeHTTP(w, req)

			// assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			mockService.AssertExpectations(t)
		})
	}
}
