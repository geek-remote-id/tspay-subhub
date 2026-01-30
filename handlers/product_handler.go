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

// GetProducts godoc
// @Summary      Get all products
// @Description  Get a list of all active products
// @Tags         product
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
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

// GetProductByID godoc
// @Summary      Get a product by ID
// @Description  Get a product by its ID
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
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

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Create a new product with the provided details
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Product Data"
// @Success      201      {object}  utils.Response
// @Failure      400      {object}  utils.Response
// @Failure      500      {object}  utils.Response
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

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update a product by ID
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id       path      int             true  "Product ID"
// @Param        product  body      models.Product  true  "Product Data"
// @Success      200      {object}  utils.Response
// @Failure      400      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Failure      500      {object}  utils.Response
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

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Soft delete a product by ID
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
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
