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
		wantErr     error
	}{
		{"valid product", "Bayam Segar", "Sayuran", 5000, nil},
		{"valid snack", "Chitato", "Snack", 12000, nil},
		{"empty name", "", "Sayuran", 5000, ErrEmptyProductName},
		{"invalid type", "Bayam", "Minuman", 5000, ErrInvalidProductType},
		{"empty type", "Bayam", "", 5000, ErrInvalidProductType},
		{"zero price", "Bayam", "Sayuran", 0, ErrInvalidPrice},
		{"negative price", "Bayam", "Sayuran", -1000, ErrInvalidPrice},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct(tt.productName, tt.productType, tt.price, "", 0)

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
