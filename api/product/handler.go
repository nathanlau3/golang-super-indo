package product

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"super-indo-api/internal/product/domain"
	"super-indo-api/internal/product/port"
	"super-indo-api/pkg/common"

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
		c.JSON(http.StatusBadRequest, common.Error(http.StatusBadRequest, "format JSON tidak valid: "+err.Error()))
		return
	}

	product, err := domain.NewProduct(req.Name, req.Type, req.Price, req.Description, req.Stock)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.createProduct.Execute(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(http.StatusInternalServerError, "gagal menyimpan produk"))
		return
	}

	c.JSON(http.StatusCreated, common.Success(http.StatusCreated, "produk berhasil ditambahkan", toProductResponse(product)))
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

	allowedFields := map[string]string{
		"name":        "name",
		"type":        "type",
		"description": "description",
	}

	filter := common.Filter{
		Search:       c.Query("search"),
		SearchFields: []string{"name", "description"},
		SortBy:       c.Query("sort_by"),
		Order:        c.Query("order"),
		Page:         page,
		Limit:        limit,
		SortAllowed: map[string]string{
			"name":  "name",
			"price": "price",
			"date":  "created_at",
			"stock": "stock",
			"type":  "type",
		},
		DefaultSort: "created_at",
	}
	columns := c.QueryArray("filter[]")
	values := c.QueryArray("filter_search[]")
	if len(columns) == 0 {
		columns = c.QueryArray("filter")
		values = c.QueryArray("filter_search")
	}
	filter.LoadFields(columns, values, allowedFields)

	products, total, err := h.getProducts.Execute(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error(http.StatusInternalServerError, "gagal mengambil data produk"))
		return
	}

	if products == nil {
		products = []domain.Product{}
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, common.Paginated(http.StatusOK, "success", toProductListResponse(products), common.Meta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}))
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error(http.StatusBadRequest, "id produk tidak valid"))
		return
	}

	product, err := h.getProductByID.Execute(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, common.Error(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error(http.StatusInternalServerError, "gagal mengambil data produk"))
		return
	}

	c.JSON(http.StatusOK, common.Success(http.StatusOK, "success", toProductResponse(product)))
}
