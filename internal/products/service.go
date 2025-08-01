package products

import (
	"context"
	"errors"
	"fmt"

	"go001/pkg/api" // Generated API types
)

// Service defines the application-level logic for product operations.
type Service struct {
	repo ProductRepository // Now depends on the repository interface
}

// NewService creates a new product service.
func NewService(repo ProductRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// FetchAllProducts retrieves a list of products.
func (s *Service) FetchAllProducts(ctx context.Context, limit, page int) ([]api.Product, error) {
	productsDB, err := s.repo.FindAll(ctx, limit, page)
	if err != nil {
		return nil, fmt.Errorf("service: failed to fetch all products from repo: %w", err)
	}
	return ToAPIProducts(productsDB), nil // Map DB models to API models
}

// GetProductByID retrieves a single product by its ID.
func (s *Service) GetProductByID(ctx context.Context, id int32) (*api.Product, error) {
	productDB, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, err) { // Propagate custom domain error
			return nil, err
		}
		return nil, fmt.Errorf("service: failed to get product from repo by ID %v: %w", id, err)
	}
	apiProduct := ToAPIProduct(productDB) // Map DB model to API model
	return &apiProduct, nil
}

// CreateProduct creates a new product
func (s *Service) CreateProduct(ctx context.Context, product api.Product) (*api.Product, error) {
	productDB, err := s.repo.Create(ctx, FromAPIProduct(product))
	if err != nil {
		return nil, fmt.Errorf("service: failed to create new product from repo: %w", err)
	}
	apiProduct := ToAPIProduct(productDB) // Map DB model to API model
	return &apiProduct, nil
}

// DeleteProduct deletes a product by ID
func (s *Service) DeleteProduct(ctx context.Context, id int32) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service: failed to delete product from repo: %w", err)
	}
	return nil
}

// UpdateProduct updates an existing product
func (s *Service) UpdateProduct(ctx context.Context, id int32, product api.Product) (*api.Product, error) {
	productDB, err := s.repo.Update(ctx, id, FromAPIProduct(product))
	if err != nil {
		return nil, fmt.Errorf("service: failed to update product from repo: %w", err)
	}
	apiProduct := ToAPIProduct(productDB) // Map DB model to API model
	return &apiProduct, nil
}
