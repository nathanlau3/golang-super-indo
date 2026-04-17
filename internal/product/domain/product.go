package domain

import (
	"errors"
	"time"
)

var (
	ErrProductNotFound    = errors.New("produk tidak ditemukan")
	ErrInvalidProductType = errors.New("tipe produk tidak valid, gunakan: Sayuran, Protein, Buah, atau Snack")
	ErrEmptyProductName   = errors.New("nama produk wajib diisi")
	ErrInvalidPrice       = errors.New("harga produk harus lebih dari 0")
)

type ProductType string

const (
	Sayuran ProductType = "Sayuran"
	Protein ProductType = "Protein"
	Buah    ProductType = "Buah"
	Snack   ProductType = "Snack"
)

var validTypes = map[ProductType]bool{
	Sayuran: true,
	Protein: true,
	Buah:    true,
	Snack:   true,
}

func (pt ProductType) IsValid() bool {
	return validTypes[pt]
}

func (pt ProductType) String() string {
	return string(pt)
}

type Product struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Type        ProductType `json:"type"`
	Price       float64     `json:"price"`
	Description string      `json:"description"`
	Stock       int         `json:"stock"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func NewProduct(name, productType string, price float64, description string, stock int) (*Product, error) {
	if name == "" {
		return nil, ErrEmptyProductName
	}

	pt := ProductType(productType)
	if !pt.IsValid() {
		return nil, ErrInvalidProductType
	}

	if price <= 0 {
		return nil, ErrInvalidPrice
	}

	return &Product{
		Name:        name,
		Type:        pt,
		Price:       price,
		Description: description,
		Stock:       stock,
	}, nil
}

type ProductFilter struct {
	Search string
	Type   ProductType
	SortBy string
	Order  string
	Page   int
	Limit  int
}

