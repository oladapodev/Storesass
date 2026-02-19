package domain

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-"`
	Role      string    `json:"role" gorm:"default:customer"`
}

type Store struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Products    []Product `json:"products,omitempty" gorm:"foreignKey:StoreID"`
}

type Product struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	StoreID     uint      `json:"store_id" gorm:"not null"`
	Store       *Store    `json:"store,omitempty" gorm:"foreignKey:StoreID"`
}

type Order struct {
	ID         uint        `json:"id" gorm:"primarykey"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	UserID     uint        `json:"user_id"`
	StoreID    uint        `json:"store_id"`
	Status     string      `json:"status" gorm:"default:pending"`
	TotalPrice float64     `json:"total_price"`
	Items      []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primarykey"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
