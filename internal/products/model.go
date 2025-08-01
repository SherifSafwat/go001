package products

import (
	"time"
)

// ProductDB represents a product in the database.
type ProductDB struct {
	ID          int32     `gorm:"type:int;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	Price       float64   `gorm:"type:numeric(10,2);not null" json:"price"`
	Category    string    `gorm:"type:varchar(100);not null" json:"category"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	//DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}

// TableName specifies the table name for GORM
func (ProductDB) TableName() string {
	return "products" // Ensures your table is named 'products'
}
