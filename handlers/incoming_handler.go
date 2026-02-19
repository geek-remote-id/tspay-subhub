package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/geek-remote-id/tspay-subhub/services"
	"github.com/geek-remote-id/tspay-subhub/utils"
)

// GenerateDepositCallbackHandler handles the deposit callback from tspay
// @Summary      Deposit Callback Handler
// @Description  Handles incoming deposit callback requests from Tspay
// @Tags         incoming
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response
// @Router       /incoming/deposit_callback [post]
func GenerateDepositCallbackHandler() http.HandlerFunc {
	tspaySvc := services.NewTspayService()
	incomingSvc := services.NewIncomingService(tspaySvc)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.Response{
				Status:  "error",
				Message: "Method not allowed",
			})
			return
		}

		// Handle webhook
		signature := r.Header.Get("X-Webhook-Signature")
		timestamp := r.Header.Get("X-Webhook-Timestamp")

		log.Printf("Incoming ctl @ deposit_callback signature = %s", signature)
		log.Printf("Incoming ctl @ deposit_callback timestamp = %s", timestamp)

		// Read body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			utils.WriteJSON(w, http.StatusBadRequest, utils.Response{
				Status:  "error",
				Message: "Could not read request body",
			})
			return
		}
		defer r.Body.Close()

		// Log request for debugging
		utils.LogToFile("deposit_callback", r.Header, body)

		// Core Logic processed in Service
		if err := incomingSvc.ProcessDepositCallback(body, signature, timestamp); err != nil {
			log.Printf("Error processing deposit callback: %v", err)
			// Decide the response based on the error
			status := http.StatusInternalServerError
			if err.Error() == "invalid signature" {
				status = http.StatusUnauthorized
			}

			utils.WriteJSON(w, status, utils.Response{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.Response{
			Status:  "success",
			Message: "Deposit Callback received",
		})
	}
}
