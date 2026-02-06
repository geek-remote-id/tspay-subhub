package handlers

import (
	"encoding/json"
	"net/http"

	"kasir-api/models"
	"kasir-api/services"
	"kasir-api/utils"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// multiple item apa aja, quantity nya
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Checkout godoc
// @Summary      Process checkout
// @Description  Create a new transaction by processing checkout items
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        checkout  body      models.CheckoutRequest  true  "Checkout Data"
// @Success      200       {object}  utils.Response
// @Failure      400       {object}  utils.Response
// @Failure      500       {object}  utils.Response
// @Router       /checkout [post]
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{
			Status:  "failed",
			Message: "Invalid request body",
		})
		return
	}

	transaction, err := h.service.Checkout(req.Items, false)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Response{
			Status:  "failed",
			Message: "Failed to process checkout: " + err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Status:  "success",
		Message: "Transaction created successfully",
		Data:    transaction,
	})
}
