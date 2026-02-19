package handlers

import (
	"fmt"
	"io"
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
			utils.LogEvent("errors", "Error reading body: "+err.Error())
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
			utils.LogEvent("errors", "Method "+r.Method+" not allowed")
			utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.Response{
				Status:  "error",
				Message: "Method not allowed",
			})
			return
		}

		// 4. Handle webhook headers
		signature := r.Header.Get("X-Webhook-Signature")
		timestamp := r.Header.Get("X-Webhook-Timestamp")
		utils.LogEvent("events", fmt.Sprintf("Incoming signature = %s, timestamp = %s", signature, timestamp))

		// 5. Core Logic processed in Service
		if err := incomingSvc.ProcessDepositCallback(body, signature, timestamp); err != nil {
			utils.LogEvent("errors", "Error processing deposit callback: "+err.Error())
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
