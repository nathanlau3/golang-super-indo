package domain

import (
	"errors"
	"testing"
)

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name        string
		productName string
		productType string
		price       float64
		description string
		stock       int
		wantErr     error
	}{
		{"valid product", "Bayam Segar", "Sayuran", 5000, "Bayam hijau segar", 100, nil},
		{"valid snack", "Chitato", "Snack", 12000, "", 0, nil},
		{"valid with all fields", "Salmon", "Protein", 89000, "Salmon fillet Norway", 20, nil},
		{"empty name", "", "Sayuran", 5000, "", 0, ErrEmptyProductName},
		{"invalid type", "Bayam", "Minuman", 5000, "", 0, ErrInvalidProductType},
		{"empty type", "Bayam", "", 5000, "", 0, ErrInvalidProductType},
		{"zero price", "Bayam", "Sayuran", 0, "", 0, ErrInvalidPrice},
		{"negative price", "Bayam", "Sayuran", -1000, "", 0, ErrInvalidPrice},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct(tt.productName, tt.productType, tt.price, tt.description, tt.stock)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewProduct() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && product == nil {
				t.Error("expected product to be non-nil")
			}

			if tt.wantErr == nil && product != nil {
				if product.Name != tt.productName {
					t.Errorf("Name = %q, want %q", product.Name, tt.productName)
				}
				if product.Type.String() != tt.productType {
					t.Errorf("Type = %q, want %q", product.Type, tt.productType)
				}
				if product.Price != tt.price {
					t.Errorf("Price = %v, want %v", product.Price, tt.price)
				}
				if product.Description != tt.description {
					t.Errorf("Description = %q, want %q", product.Description, tt.description)
				}
				if product.Stock != tt.stock {
					t.Errorf("Stock = %d, want %d", product.Stock, tt.stock)
				}
				if product.DeletedAt != nil {
					t.Error("DeletedAt should be nil for new product")
				}
			}
		})
	}
}

func TestProductType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		pt       ProductType
		expected bool
	}{
		{"Sayuran valid", Sayuran, true},
		{"Protein valid", Protein, true},
		{"Buah valid", Buah, true},
		{"Snack valid", Snack, true},
		{"Minuman invalid", ProductType("Minuman"), false},
		{"empty invalid", ProductType(""), false},
		{"lowercase invalid", ProductType("sayuran"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pt.IsValid(); got != tt.expected {
				t.Errorf("ProductType(%q).IsValid() = %v, want %v", tt.pt, got, tt.expected)
			}
		})
	}
}
