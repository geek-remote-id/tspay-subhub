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
		// 1. Read body first so we can log it (Body can only be read once)
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

		// 2. Log EVERYTHING immediately to file
		utils.LogToFile("deposit_callback", r.Header, body)

		// 3. Method check
		if r.Method != http.MethodPost {
			log.Printf("Method %s not allowed", r.Method)
			utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.Response{
				Status:  "error",
				Message: "Method not allowed",
			})
			return
		}

		// 4. Handle webhook headers
		signature := r.Header.Get("X-Webhook-Signature")
		timestamp := r.Header.Get("X-Webhook-Timestamp")
		log.Printf("Incoming signature = %s, timestamp = %s", signature, timestamp)

		// 5. Core Logic processed in Service
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
