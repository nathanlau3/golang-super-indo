package postgres

import (
	"database/sql"
	"log"
)

type seedProduct struct {
	Name        string
	Type        string
	Price       float64
	Description string
	Stock       int
}

func SeedProducts(db *sql.DB) {
	products := []seedProduct{
		{"Bayam Segar", "Sayuran", 5500, "Bayam hijau segar 250g", 100},
		{"Kangkung", "Sayuran", 3500, "Kangkung segar 200g", 120},
		{"Brokoli Import", "Sayuran", 18000, "Brokoli import per 500g", 45},
		{"Wortel Lokal", "Sayuran", 9000, "Wortel lokal per kg", 80},
		{"Selada Romaine", "Sayuran", 7500, "Selada romaine segar 150g", 55},

		{"Dada Ayam Fillet", "Protein", 38000, "Dada ayam tanpa tulang per kg", 40},
		{"Ikan Salmon Fillet", "Protein", 89000, "Salmon fillet Norway 200g", 20},
		{"Telur Ayam Kampung", "Protein", 32000, "Telur ayam kampung 10 butir", 65},
		{"Daging Sapi Has Dalam", "Protein", 135000, "Tenderloin sapi lokal per kg", 15},
		{"Udang Vaname", "Protein", 62000, "Udang vaname segar size 40 per kg", 30},

		{"Apel Fuji", "Buah", 28000, "Apel Fuji import per kg", 70},
		{"Pisang Cavendish", "Buah", 19500, "Pisang cavendish per sisir", 90},
		{"Jeruk Mandarin", "Buah", 32000, "Jeruk mandarin per kg", 50},
		{"Semangka Merah", "Buah", 24000, "Semangka merah tanpa biji per buah", 25},
		{"Mangga Harum Manis", "Buah", 22000, "Mangga harum manis per kg", 60},

		{"Chitato Original 68g", "Snack", 11500, "Keripik kentang rasa original", 200},
		{"Pocky Chocolate", "Snack", 10000, "Pocky stick rasa coklat", 180},
		{"Oreo Original 133g", "Snack", 9500, "Biskuit sandwich oreo original", 210},
		{"Tango Wafer Coklat", "Snack", 8500, "Wafer coklat Tango 176g", 150},
		{"Lays Rumput Laut 55g", "Snack", 10500, "Keripik kentang rasa rumput laut", 175},
	}

	for _, p := range products {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM products WHERE name = $1", p.Name).Scan(&count)
		if err != nil {
			log.Printf("gagal cek produk %s: %v", p.Name, err)
			continue
		}
		if count > 0 {
			log.Printf("skip: %s (sudah ada)", p.Name)
			continue
		}

		_, err = db.Exec(
			`INSERT INTO products (name, type, price, description, stock, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
			p.Name, p.Type, p.Price, p.Description, p.Stock,
		)
		if err != nil {
			log.Printf("gagal seed %s: %v", p.Name, err)
		} else {
			log.Printf("seeded: %s", p.Name)
		}
	}
}
