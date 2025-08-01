package products

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"go001/pkg/api" // Generated API types
)

// ProductRepository defines the interface for product data operations.
type ProductRepository interface {
	FindAll(ctx context.Context, limit, page int) ([]ProductDB, error)
	FindByID(ctx context.Context, id int32) (*ProductDB, error)
	Create(ctx context.Context, ProductDB ProductDB) (*ProductDB, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, id int32, product ProductDB) (*ProductDB, error)
}

// GORMProductRepository implements ProductRepository using GORM.
type GORMProductRepository struct {
	db *gorm.DB
}

// NewGORMProductRepository creates a new GORMProductRepository.
func NewGORMProductRepository(db *gorm.DB) *GORMProductRepository {
	return &GORMProductRepository{db: db}
}

// FindAll retrieves a list of products from the database.
func (r *GORMProductRepository) FindAll(ctx context.Context, limit, page int) ([]ProductDB, error) {
	var products []ProductDB

	if page < 0 {
		return nil, fmt.Errorf("failed to fetch products: wrong page number")
	}

	offset := (page - 1) * limit
	query := r.db.WithContext(ctx).Limit(limit).Offset(offset)
	result := query.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", result.Error)
	}
	return products, nil
}

// FindByID retrieves a single product by its ID from the database.
func (r *GORMProductRepository) FindByID(ctx context.Context, id int32) (*ProductDB, error) {
	var product ProductDB
	result := r.db.WithContext(ctx).First(&product, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error // Use our custom domain error
		}
		return nil, fmt.Errorf("failed to find product by ID %v: %w", id, result.Error)
	}
	return &product, nil
}

// Create creates a new product
func (r *GORMProductRepository) Create(ctx context.Context, product ProductDB) (*ProductDB, error) {
	result := r.db.WithContext(ctx).Create(&product)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create product: %w", result.Error)
	}
	return &product, nil
}

// Delete deletes a product by ID
func (r *GORMProductRepository) Delete(ctx context.Context, id int32) error {
	result := r.db.WithContext(ctx).Delete(&ProductDB{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Update updates an existing product
func (r *GORMProductRepository) Update(ctx context.Context, id int32, product ProductDB) (*ProductDB, error) {
	result := r.db.WithContext(ctx).Model(&ProductDB{}).Where("id = ?", id).Updates(product)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var updatedProduct ProductDB
	if err := r.db.WithContext(ctx).First(&updatedProduct, id).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve updated product: %w", err)
	}
	return &updatedProduct, nil
}

// --- Mapper functions between API and DB models ---

// ToAPIProduct converts a ProductDB model to an api.Product model.
func ToAPIProduct(pDB *ProductDB) api.Product {
	if pDB == nil {
		return api.Product{} // Return empty or handle nil appropriately
	}
	return api.Product{
		Id:          pDB.ID,
		Name:        pDB.Name,
		Description: &pDB.Description,
		Price:       pDB.Price,
		Category:    pDB.Category,
	}
}

// ToAPIProducts converts a slice of ProductDB models to a slice of api.Product models.
func ToAPIProducts(pDBs []ProductDB) []api.Product {
	apiProducts := make([]api.Product, len(pDBs))
	for i, pDB := range pDBs {
		apiProducts[i] = ToAPIProduct(&pDB)
	}
	return apiProducts
}

// FromAPIProduct converts an api.Product model to a ProductDB model.
// Useful if you have a POST/PUT endpoint to convert incoming API data to DB data.
func FromAPIProduct(apiP api.Product) ProductDB {

	return ProductDB{
		ID:          apiP.Id,
		Name:        apiP.Name,
		Description: *apiP.Description,
		Price:       apiP.Price,
		Category:    apiP.Category,
	}
}
