package main

import (
	"log"

	"github.com/yourusername/storefront-saas-go/internal/config"
	"github.com/yourusername/storefront-saas-go/internal/db"
	"github.com/yourusername/storefront-saas-go/internal/domain"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	stores := []domain.Store{
		{Name: "Demo Store One", Slug: "demo-store-one", Description: "A demo storefront for testing", IsActive: true},
		{Name: "Sample Marketplace", Slug: "sample-marketplace", Description: "A sample marketplace", IsActive: true},
	}

	for i := range stores {
		if err := database.FirstOrCreate(&stores[i], domain.Store{Slug: stores[i].Slug}).Error; err != nil {
			log.Printf("failed to seed store %s: %v", stores[i].Slug, err)
		} else {
			log.Printf("seeded store: %s (id=%d)", stores[i].Name, stores[i].ID)
		}
	}

	products := []domain.Product{
		{Name: "Basic Widget", Description: "A basic widget", Price: 9.99, Stock: 100, IsActive: true, StoreID: stores[0].ID},
		{Name: "Premium Tool", Description: "A premium tool", Price: 49.99, Stock: 50, IsActive: true, StoreID: stores[0].ID},
		{Name: "Standard Package", Description: "A standard package", Price: 29.99, Stock: 75, IsActive: true, StoreID: stores[0].ID},
		{Name: "Essential Kit", Description: "An essential kit", Price: 19.99, Stock: 200, IsActive: true, StoreID: stores[1].ID},
		{Name: "Advanced Module", Description: "An advanced module", Price: 79.99, Stock: 30, IsActive: true, StoreID: stores[1].ID},
		{Name: "Professional Bundle", Description: "A professional bundle", Price: 99.99, Stock: 20, IsActive: true, StoreID: stores[1].ID},
	}

	for i := range products {
		if err := database.FirstOrCreate(&products[i], domain.Product{Name: products[i].Name, StoreID: products[i].StoreID}).Error; err != nil {
			log.Printf("failed to seed product %s: %v", products[i].Name, err)
		} else {
			log.Printf("seeded product: %s (id=%d)", products[i].Name, products[i].ID)
		}
	}

	// Seed sample order (2 × 9.99 + 1 × 49.99 = 69.97)
	order := domain.Order{
		UserID:     1,
		StoreID:    stores[0].ID,
		Status:     "completed",
		TotalPrice: 69.97,
		Items: []domain.OrderItem{
			{ProductID: products[0].ID, Quantity: 2, Price: 9.99},
			{ProductID: products[1].ID, Quantity: 1, Price: 49.99},
		},
	}

	if err := database.Create(&order).Error; err != nil {
		log.Printf("failed to seed order: %v", err)
	} else {
		log.Printf("seeded order id=%d", order.ID)
	}

	log.Println("Seeding complete!")
}
