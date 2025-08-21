package order

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GkadyrG/L0/backend/internal/apperr"
	"github.com/GkadyrG/L0/backend/internal/model"
	"github.com/GkadyrG/L0/backend/internal/repository/mocks"
	"github.com/GkadyrG/L0/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestHandler(repo *mocks.OrderRepository) *Handler {
	uc := usecase.New(repo)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return New(uc, logger)
}

func TestHandler_GetByID(t *testing.T) {
	type testCase struct {
		name       string
		id         string
		mockSetup  func(r *mocks.OrderRepository)
		wantCode   int
		assertBody func(t *testing.T, body []byte)
	}

	now := time.Now().UTC().Truncate(time.Second)
	expectedOrder := &model.OrderResponse{
		OrderUID:    "order-1",
		TrackNumber: "TRK123",
		CustomerID:  "cust-1",
		DateCreated: now,
		Delivery: model.DeliveryResponse{
			Name:    "John Doe",
			Phone:   "+1234567",
			City:    "City",
			Address: "Street 1",
			Email:   "john@example.com",
		},
		Payment: model.PaymentResponse{
			Transaction: "tx-1",
			Currency:    "USD",
			Amount:      100,
		},
		Items: []model.ItemResponse{{
			Name:   "Item A",
			Price:  100,
			Brand:  "BrandX",
			Status: 1,
		}},
	}

	tests := []testCase{
		{
			name: "success",
			id:   "order-1",
			mockSetup: func(r *mocks.OrderRepository) {
				r.On("GetByID", mock.Anything, "order-1").Return(expectedOrder, nil)
			},
			wantCode: http.StatusOK,
			assertBody: func(t *testing.T, body []byte) {
				var got model.OrderResponse
				assert.NoError(t, json.Unmarshal(body, &got))
				exp := *expectedOrder
				exp.DateCreated = got.DateCreated
				assert.Equal(t, exp, got)
			},
		},
		{
			name: "not found",
			id:   "missing",
			mockSetup: func(r *mocks.OrderRepository) {
				r.On("GetByID", mock.Anything, "missing").Return((*model.OrderResponse)(nil), apperr.ErrNotFound)
			},
			wantCode: http.StatusNotFound,
			assertBody: func(t *testing.T, body []byte) {
				var m map[string]string
				assert.NoError(t, json.Unmarshal(body, &m))
				assert.Equal(t, "order not found", m["error"])
			},
		},
		{
			name: "internal error",
			id:   "boom",
			mockSetup: func(r *mocks.OrderRepository) {
				r.On("GetByID", mock.Anything, "boom").Return((*model.OrderResponse)(nil), assert.AnError)
			},
			wantCode: http.StatusInternalServerError,
			assertBody: func(t *testing.T, body []byte) {
				var m map[string]string
				assert.NoError(t, json.Unmarshal(body, &m))
				assert.Equal(t, "internal server error", m["error"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewOrderRepository(t)
			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}
			h := newTestHandler(repo)

			router := chi.NewRouter()
			router.Get("/api/order/{id}", h.GetByID())

			req := httptest.NewRequest(http.MethodGet, "/api/order/"+tc.id, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
			if tc.assertBody != nil {
				tc.assertBody(t, rec.Body.Bytes())
			}
		})
	}
}

func TestHandler_GetAll(t *testing.T) {
	type testCase struct {
		name       string
		mockSetup  func(r *mocks.OrderRepository)
		wantCode   int
		assertBody func(t *testing.T, body []byte)
	}

	now := time.Now().UTC().Truncate(time.Second)
	preview := []*model.OrderPreview{{
		OrderUID:    "order-1",
		TrackNumber: "TRK123",
		CustomerID:  "cust-1",
		DateCreated: now,
	}}

	tests := []testCase{
		{
			name: "success",
			mockSetup: func(r *mocks.OrderRepository) {
				r.On("GetAll", mock.Anything).Return(preview, nil)
			},
			wantCode: http.StatusOK,
			assertBody: func(t *testing.T, body []byte) {
				var got []*model.OrderPreview
				assert.NoError(t, json.Unmarshal(body, &got))
				if len(got) > 0 && len(preview) > 0 {
					exp := *preview[0]
					exp.DateCreated = got[0].DateCreated
					assert.Equal(t, exp, *got[0])
				} else {
					assert.Equal(t, len(preview), len(got))
				}
			},
		},
		{
			name: "not found",
			mockSetup: func(r *mocks.OrderRepository) {
				r.On("GetAll", mock.Anything).Return(([]*model.OrderPreview)(nil), apperr.ErrNotFound)
			},
			wantCode: http.StatusNotFound,
			assertBody: func(t *testing.T, body []byte) {
				var m map[string]string
				assert.NoError(t, json.Unmarshal(body, &m))
				assert.Equal(t, "orders preview not found", m["error"])
			},
		},
		{
			name: "internal error",
			mockSetup: func(r *mocks.OrderRepository) {
				r.On("GetAll", mock.Anything).Return(([]*model.OrderPreview)(nil), assert.AnError)
			},
			wantCode: http.StatusInternalServerError,
			assertBody: func(t *testing.T, body []byte) {
				var m map[string]string
				assert.NoError(t, json.Unmarshal(body, &m))
				assert.Equal(t, "internal server error", m["error"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewOrderRepository(t)
			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}
			h := newTestHandler(repo)

			router := chi.NewRouter()
			router.Get("/api/orders", h.GetAll())

			req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
			if tc.assertBody != nil {
				tc.assertBody(t, rec.Body.Bytes())
			}
		})
	}
}
