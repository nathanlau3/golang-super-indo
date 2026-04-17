package product

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"super-indo-api/internal/product/domain"
	"super-indo-api/internal/product/port"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	createProduct  port.CreateProductUseCase
	getProducts    port.GetProductsUseCase
	getProductByID port.GetProductByIDUseCase
}

func NewProductHandler(
	createProduct port.CreateProductUseCase,
	getProducts port.GetProductsUseCase,
	getProductByID port.GetProductByIDUseCase,
) *ProductHandler {
	return &ProductHandler{
		createProduct:  createProduct,
		getProducts:    getProducts,
		getProductByID: getProductByID,
	}
}

func (h *ProductHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.GetProducts)
	rg.GET("/:id", h.GetProductByID)
	rg.POST("", h.CreateProduct)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  http.StatusBadRequest,
			Message: "format JSON tidak valid: " + err.Error(),
		})
		return
	}

	product, err := domain.NewProduct(req.Name, req.Type, req.Price, req.Description, req.Stock)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if err := h.createProduct.Execute(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  http.StatusInternalServerError,
			Message: "gagal menyimpan produk",
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: "produk berhasil ditambahkan",
		Data:    toProductResponse(product),
	})
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	typeStr := c.Query("type")
	var productType domain.ProductType
	if typeStr != "" {
		productType = domain.ProductType(typeStr)
		if !productType.IsValid() {
			c.JSON(http.StatusBadRequest, Response{
				Status:  http.StatusBadRequest,
				Message: "tipe produk tidak valid",
			})
			return
		}
	}

	filter := domain.ProductFilter{
		Search: c.Query("search"),
		Type:   productType,
		SortBy: c.Query("sort_by"),
		Order:  c.Query("order"),
		Page:   page,
		Limit:  limit,
	}

	products, total, err := h.getProducts.Execute(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  http.StatusInternalServerError,
			Message: "gagal mengambil data produk",
		})
		return
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, PaginatedResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    toProductListResponse(products),
		Meta: Meta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  http.StatusBadRequest,
			Message: "id produk tidak valid",
		})
		return
	}

	product, err := h.getProductByID.Execute(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, Response{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Status:  http.StatusInternalServerError,
			Message: "gagal mengambil data produk",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    toProductResponse(product),
	})
}
