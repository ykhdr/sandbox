package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
}

type Product struct {
	ID       uint    `gorm:"primaryKey"`
	Name     string  `gorm:"type:varchar(100);not null"`
	Price    float64 `gorm:"type:decimal(10,2);not null"`
	Category string  `gorm:"type:varchar(50);not null"`
}

type Order struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	ProductID  uint      `gorm:"not null"`
	Quantity   int       `gorm:"not null"`
	TotalPrice float64   `gorm:"type:decimal(10,2);not null"`
	Status     string    `gorm:"type:varchar(20);not null"`
	CreatedAt  time.Time `gorm:"not null"`

	User    User    `gorm:"foreignKey:UserID"`
	Product Product `gorm:"foreignKey:ProductID"`
}

// views

type OrderDetail struct {
	OrderID     uint      `gorm:"column:order_id"`
	UserName    string    `gorm:"column:user_name"`
	ProductName string    `gorm:"column:product_name"`
	Quantity    int       `gorm:"column:quantity"`
	TotalPrice  float64   `gorm:"column:total_price"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (OrderDetail) TableName() string {
	return "order_details"
}

type MonthlySalesSummary struct {
	Month        time.Time `gorm:"column:month"`
	TotalOrders  int       `gorm:"column:total_orders"`
	TotalRevenue float64   `gorm:"column:total_revenue"`
	AvgCheck     float64   `gorm:"column:avg_check"`
}

func (MonthlySalesSummary) TableName() string {
	return "monthly_sales_summary"
}
