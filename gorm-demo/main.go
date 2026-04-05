package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gorm-demo/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=demo password=demo dbname=gorm_demo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	fmt.Println("connected to database")

	runMigrations(db)

	seedData(db)
	demoCRUD(db)
	demoViews(db)
}

func runMigrations(db *gorm.DB) {
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatal("failed to read migrations:", err)
	}
	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("failed to read %s: %v", file, err)
		}
		if err := db.Exec(string(content)).Error; err != nil {
			log.Fatalf("failed to execute %s: %v", file, err)
		}
		fmt.Printf("migration applied: %s\n", file)
	}
}

func seedData(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		fmt.Println("data already seeded, skipping")
		return
	}

	users := []models.User{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
		{Name: "Charlie", Email: "charlie@example.com"},
	}
	db.Create(&users)

	products := []models.Product{
		{Name: "Laptop", Price: 1200.00, Category: "electronics"},
		{Name: "Keyboard", Price: 75.00, Category: "electronics"},
		{Name: "Go Book", Price: 35.00, Category: "books"},
		{Name: "Coffee Mug", Price: 12.00, Category: "home"},
		{Name: "Monitor", Price: 450.00, Category: "electronics"},
	}
	db.Create(&products)

	orders := []models.Order{
		{UserID: 1, ProductID: 1, Quantity: 1, TotalPrice: 1200.00, Status: "completed", CreatedAt: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)},
		{UserID: 1, ProductID: 3, Quantity: 2, TotalPrice: 70.00, Status: "completed", CreatedAt: time.Date(2025, 1, 20, 14, 0, 0, 0, time.UTC)},
		{UserID: 2, ProductID: 2, Quantity: 1, TotalPrice: 75.00, Status: "completed", CreatedAt: time.Date(2025, 2, 5, 9, 0, 0, 0, time.UTC)},
		{UserID: 2, ProductID: 5, Quantity: 1, TotalPrice: 450.00, Status: "shipped", CreatedAt: time.Date(2025, 2, 18, 16, 0, 0, 0, time.UTC)},
		{UserID: 3, ProductID: 4, Quantity: 3, TotalPrice: 36.00, Status: "completed", CreatedAt: time.Date(2025, 3, 1, 11, 0, 0, 0, time.UTC)},
		{UserID: 1, ProductID: 5, Quantity: 2, TotalPrice: 900.00, Status: "pending", CreatedAt: time.Date(2025, 3, 10, 8, 0, 0, 0, time.UTC)},
		{UserID: 3, ProductID: 1, Quantity: 1, TotalPrice: 1200.00, Status: "completed", CreatedAt: time.Date(2025, 4, 5, 13, 0, 0, 0, time.UTC)},
		{UserID: 2, ProductID: 3, Quantity: 1, TotalPrice: 35.00, Status: "shipped", CreatedAt: time.Date(2025, 4, 22, 10, 0, 0, 0, time.UTC)},
		{UserID: 1, ProductID: 4, Quantity: 5, TotalPrice: 60.00, Status: "pending", CreatedAt: time.Date(2025, 5, 3, 15, 0, 0, 0, time.UTC)},
		{UserID: 3, ProductID: 2, Quantity: 2, TotalPrice: 150.00, Status: "completed", CreatedAt: time.Date(2025, 5, 15, 12, 0, 0, 0, time.UTC)},
	}
	db.Create(&orders)

	fmt.Printf("seeded: %d users, %d products, %d orders\n", len(users), len(products), len(orders))
}

func demoCRUD(db *gorm.DB) {
	fmt.Println("\n=== CRUD Demo ===")

	newOrder := models.Order{
		UserID: 2, ProductID: 1, Quantity: 1,
		TotalPrice: 1200.00, Status: "new",
		CreatedAt: time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC),
	}
	db.Create(&newOrder)
	fmt.Printf("created order: id=%d\n", newOrder.ID)

	var order models.Order
	db.Preload("User").Preload("Product").First(&order, "id = ? AND created_at = ?", newOrder.ID, newOrder.CreatedAt)
	fmt.Printf("read order: %s bought %s, total=%.2f\n", order.User.Name, order.Product.Name, order.TotalPrice)

	db.Model(&order).Where("created_at = ?", order.CreatedAt).Update("status", "processing")
	fmt.Printf("updated order status to: processing\n")

	db.Where("id = ? AND created_at = ?", order.ID, order.CreatedAt).Delete(&models.Order{})
	fmt.Printf("deleted order: id=%d\n", order.ID)
}

func demoViews(db *gorm.DB) {
	fmt.Println("\n=== Order Details View ===")
	var details []models.OrderDetail
	db.Find(&details)
	for _, d := range details {
		fmt.Printf("  #%d | %-8s | %-10s | qty=%d | $%.2f | %s | %s\n",
			d.OrderID, d.UserName, d.ProductName, d.Quantity, d.TotalPrice, d.Status, d.CreatedAt.Format("2006-01"))
	}

	fmt.Println("\n=== Monthly Sales Summary View ===")
	var summary []models.MonthlySalesSummary
	db.Find(&summary)
	for _, s := range summary {
		fmt.Printf("  %s | orders=%d | revenue=$%.2f | avg=$%.2f\n",
			s.Month.Format("2006-01"), s.TotalOrders, s.TotalRevenue, s.AvgCheck)
	}
}
