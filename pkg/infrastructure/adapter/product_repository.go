package adapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"super-indo-api/internal/product/domain"
	"super-indo-api/internal/product/port"

	"github.com/redis/go-redis/v9"
)

const cacheTTL = 5 * time.Minute

type productRepository struct {
	db    *sql.DB
	cache *redis.Client
}

func NewProductRepository(db *sql.DB, cache *redis.Client) port.ProductRepository {
	return &productRepository{db: db, cache: cache}
}

func (r *productRepository) Save(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (name, type, price, description, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		product.Name, product.Type.String(), product.Price, product.Description, product.Stock,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return err
	}

	r.invalidateListCache(ctx)
	return nil
}

func (r *productRepository) FindAll(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int64, error) {
	cacheKey := fmt.Sprintf("products:list:%s:%s:%s:%s:%s:%s:%d:%d",
		filter.Search, filter.Name, filter.Type, filter.Description,
		filter.SortBy, filter.Order, filter.Page, filter.Limit,
	)

	cached, err := r.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var result cachedProductList
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Products, result.Total, nil
		}
	}

	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(LOWER(name) LIKE LOWER($%d) OR LOWER(COALESCE(description, '')) LIKE LOWER($%d))", argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(name) LIKE LOWER($%d)", argIdx))
		args = append(args, "%"+filter.Name+"%")
		argIdx++
	}

	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(type) LIKE LOWER($%d)", argIdx))
		args = append(args, "%"+string(filter.Type)+"%")
		argIdx++
	}

	if filter.Description != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(COALESCE(description, '')) LIKE LOWER($%d)", argIdx))
		args = append(args, "%"+filter.Description+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int64
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count query: %w", err)
	}

	sortColumn := "created_at"
	switch filter.SortBy {
	case "price":
		sortColumn = "price"
	case "name":
		sortColumn = "name"
	case "date":
		sortColumn = "created_at"
	}

	orderDir := "DESC"
	if filter.Order == "asc" {
		orderDir = "ASC"
	}

	offset := (filter.Page - 1) * filter.Limit

	dataSQL := fmt.Sprintf(`
		SELECT id, name, type, price, COALESCE(description, ''), stock, created_at, updated_at
		FROM products %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		whereClause, sortColumn, orderDir, argIdx, argIdx+1,
	)
	args = append(args, filter.Limit, offset)

	rows, err := r.db.QueryContext(ctx, dataSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("select query: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Type, &p.Price,
			&p.Description, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan row: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	go func() {
		data, _ := json.Marshal(cachedProductList{Products: products, Total: total})
		if err := r.cache.Set(context.Background(), cacheKey, data, cacheTTL).Err(); err != nil {
			log.Printf("cache set error: %v", err)
		}
	}()

	return products, total, nil
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*domain.Product, error) {
	cacheKey := fmt.Sprintf("products:detail:%d", id)

	cached, err := r.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var product domain.Product
		if json.Unmarshal([]byte(cached), &product) == nil {
			return &product, nil
		}
	}

	query := `
		SELECT id, name, type, price, COALESCE(description, ''), stock, created_at, updated_at
		FROM products
		WHERE id = $1`

	var p domain.Product
	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Type, &p.Price,
		&p.Description, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	data, _ := json.Marshal(p)
	r.cache.Set(ctx, cacheKey, data, cacheTTL)

	return &p, nil
}

func (r *productRepository) invalidateListCache(ctx context.Context) {
	iter := r.cache.Scan(ctx, 0, "products:list:*", 100).Iterator()
	for iter.Next(ctx) {
		r.cache.Del(ctx, iter.Val())
	}
}

type cachedProductList struct {
	Products []domain.Product `json:"products"`
	Total    int64            `json:"total"`
}
