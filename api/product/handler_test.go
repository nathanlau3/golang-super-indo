package product

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"super-indo-api/internal/product/domain"

	"github.com/gin-gonic/gin"
)

// mock per use case
type mockCreateProduct struct {
	fn func(ctx context.Context, p *domain.Product) error
}

func (m *mockCreateProduct) Execute(ctx context.Context, p *domain.Product) error {
	if m.fn != nil {
		return m.fn(ctx, p)
	}
	return nil
}

type mockGetProducts struct {
	fn func(ctx context.Context, f domain.ProductFilter) ([]domain.Product, int64, error)
}

func (m *mockGetProducts) Execute(ctx context.Context, f domain.ProductFilter) ([]domain.Product, int64, error) {
	if m.fn != nil {
		return m.fn(ctx, f)
	}
	return nil, 0, nil
}

type mockGetProductByID struct {
	fn func(ctx context.Context, id uint) (*domain.Product, error)
}

func (m *mockGetProductByID) Execute(ctx context.Context, id uint) (*domain.Product, error) {
	if m.fn != nil {
		return m.fn(ctx, id)
	}
	return nil, nil
}

func setupTestRouter(h *ProductHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h.RegisterRoutes(r.Group("/product"))
	return r
}

func newTestHandler(
	create *mockCreateProduct,
	list *mockGetProducts,
	byID *mockGetProductByID,
) *ProductHandler {
	if create == nil {
		create = &mockCreateProduct{}
	}
	if list == nil {
		list = &mockGetProducts{}
	}
	if byID == nil {
		byID = &mockGetProductByID{}
	}
	return NewProductHandler(create, list, byID)
}

func TestCreateProduct_Success(t *testing.T) {
	h := newTestHandler(
		&mockCreateProduct{fn: func(ctx context.Context, p *domain.Product) error {
			p.ID = 1
			return nil
		}},
		nil, nil,
	)
	r := setupTestRouter(h)

	body := map[string]interface{}{
		"name":  "Bayam Segar",
		"type":  "Sayuran",
		"price": 5000,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestCreateProduct_InvalidType(t *testing.T) {
	h := newTestHandler(nil, nil, nil)
	r := setupTestRouter(h)

	body := map[string]interface{}{
		"name":  "Teh Botol",
		"type":  "Minuman",
		"price": 5000,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateProduct_MissingName(t *testing.T) {
	h := newTestHandler(nil, nil, nil)
	r := setupTestRouter(h)

	body := map[string]interface{}{
		"type":  "Sayuran",
		"price": 5000,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateProduct_InvalidBody(t *testing.T) {
	h := newTestHandler(nil, nil, nil)
	r := setupTestRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/product", bytes.NewBufferString("bukan json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetProducts_Success(t *testing.T) {
	products := []domain.Product{
		{ID: 1, Name: "Bayam Segar", Type: domain.Sayuran, Price: 5000},
		{ID: 2, Name: "Kangkung", Type: domain.Sayuran, Price: 3000},
	}

	h := newTestHandler(nil,
		&mockGetProducts{fn: func(ctx context.Context, f domain.ProductFilter) ([]domain.Product, int64, error) {
			return products, 2, nil
		}},
		nil,
	)
	r := setupTestRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/product?type=Sayuran&page=1&limit=10", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp PaginatedResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Meta.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Meta.Total)
	}
}

func TestGetProducts_InvalidType(t *testing.T) {
	h := newTestHandler(nil, nil, nil)
	r := setupTestRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/product?type=Minuman", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetProductByID_Success(t *testing.T) {
	h := newTestHandler(nil, nil,
		&mockGetProductByID{fn: func(ctx context.Context, id uint) (*domain.Product, error) {
			return &domain.Product{ID: 1, Name: "Bayam Segar", Type: domain.Sayuran, Price: 5000}, nil
		}},
	)
	r := setupTestRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/product/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	h := newTestHandler(nil, nil,
		&mockGetProductByID{fn: func(ctx context.Context, id uint) (*domain.Product, error) {
			return nil, domain.ErrProductNotFound
		}},
	)
	r := setupTestRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/product/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetProductByID_InvalidID(t *testing.T) {
	h := newTestHandler(nil, nil, nil)
	r := setupTestRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/product/abc", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
