package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go001/internal/products" // Your business logic service
	"go001/pkg/api"           // Generated types and interface
)

// APIHandler implements the api.ServerInterface
type APIHandler struct {
	productService *products.Service
	// Add other services or dependencies here
}

// NewAPIHandler creates a new instance of your API handler implementation
func NewAPIHandler(productService *products.Service) *APIHandler {
	return &APIHandler{
		productService: productService,
	}
}

// --- Implement the methods from api.ServerInterface ---

// GetAllProducts handles GET /products
func (h *APIHandler) GetProducts(w http.ResponseWriter, r *http.Request, params api.GetProductsParams) {
	ctx := r.Context()

	limit := 20
	if params.Limit != nil {
		limit = int(*params.Limit)
	}
	page := 0
	if params.Page != nil {
		page = int(*params.Page)
	}

	productsList, err := h.productService.FetchAllProducts(ctx, limit, page)
	if err != nil {
		// Log the error for debugging, but return a generic message to client for internal errors
		fmt.Printf("Error fetching products: %v\n", err)
		// Check for specific service errors if needed
		api.ErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(productsList); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		api.ErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// GetProductById handles GET /products/{id}
func (h *APIHandler) GetProductById(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()

	product, err := h.productService.GetProductByID(ctx, int32(id))
	if err != nil {
		fmt.Printf("Error getting product by ID: %v\n", err)
		api.ErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		api.ErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// CreateProduct creates a new product
func (h *APIHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var product api.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		api.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdProduct, err := h.productService.CreateProduct(ctx, product)
	if err != nil {
		fmt.Printf("Error creating product: %v\n", err)
		api.ErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdProduct); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		api.ErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// DeleteProduct deletes a product by ID
func (h *APIHandler) DeleteProduct(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()

	if err := h.productService.DeleteProduct(ctx, int32(id)); err != nil {
		fmt.Printf("Error deleting product: %v\n", err)
		api.ErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateProduct updates a product
func (h *APIHandler) UpdateProduct(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()

	var product api.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		api.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedProduct, err := h.productService.UpdateProduct(ctx, int32(id), product)
	if err != nil {
		fmt.Printf("Error updating product: %v\n", err)
		api.ErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedProduct); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		api.ErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
	}
}
