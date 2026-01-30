package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"kasir-api/models"
	"kasir-api/services"
	"kasir-api/utils"
)

type ProductHandler struct {
	Service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{Service: service}
}

// @Router       /product [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.Service.GetAll()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Response{
			Status:  "failed",
			Message: "Failed to fetch products: " + err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Status:  "success",
		Message: "Products retrieved successfully",
		Data:    products,
	})
}

// @Router       /product/{id} [get]
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{
			Status:  "failed",
			Message: "Invalid Product ID",
		})
		return
	}

	product, err := h.Service.GetByID(id)
	if err == sql.ErrNoRows {
		utils.WriteJSON(w, http.StatusNotFound, utils.Response{
			Status:  "failed",
			Message: "Product not found",
		})
		return
	}

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Response{
			Status:  "failed",
			Message: "Failed to fetch product: " + err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Status:  "success",
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

// @Router       /product [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productReq models.Product
	err := json.NewDecoder(r.Body).Decode(&productReq)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{
			Status:  "failed",
			Message: "Invalid request body",
		})
		return
	}

	product, err := h.Service.Create(productReq)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Response{
			Status:  "failed",
			Message: "Failed to save product: " + err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Response{
		Status:  "success",
		Message: "Product created successfully",
		Data:    product,
	})
}

// @Router       /product/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{
			Status:  "failed",
			Message: "Invalid Product ID",
		})
		return
	}

	var updateReq models.Product
	err = json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{
			Status:  "failed",
			Message: "Invalid request body",
		})
		return
	}

	existingProduct, err := h.Service.GetByID(id)
	if err == sql.ErrNoRows {
		utils.WriteJSON(w, http.StatusNotFound, utils.Response{
			Status:  "failed",
			Message: "Product not found",
		})
		return
	}
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Response{
			Status:  "failed",
			Message: "Failed to fetch product: " + err.Error(),
		})
		return
	}

	if updateReq.Name != "" {
		existingProduct.Name = updateReq.Name
	}
	if updateReq.Price != 0 {
		existingProduct.Price = updateReq.Price
	}
	if updateReq.Stock != 0 {
		existingProduct.Stock = updateReq.Stock
	}
	if updateReq.CategoryID != 0 {
		existingProduct.CategoryID = updateReq.CategoryID
	}

	updatedProduct, err := h.Service.Update(existingProduct)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Response{
			Status:  "failed",
			Message: "Failed to update product: " + err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Status:  "success",
		Message: "Product updated successfully",
		Data:    updatedProduct,
	})
}

// @Router       /product/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{
			Status:  "failed",
			Message: "Invalid Product ID",
		})
		return
	}

	err = h.Service.Delete(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Response{
			Status:  "failed",
			Message: "Failed to delete product: " + err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Status:  "success",
		Message: "Product deleted successfully",
	})
}
