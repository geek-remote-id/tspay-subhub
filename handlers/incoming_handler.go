package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/geek-remote-id/tspay-subhub/models"
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
	incomingSvc := services.NewIncomingService()
	tspaySvc := services.NewTspayService()

	return func(w http.ResponseWriter, r *http.Request) {
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

		log.Printf("Incoming ctl @ deposit_callback data = %s", string(body))

		var data models.DepositCallback
		if err := json.Unmarshal(body, &data); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			// In Go, if it's not JSON, we don't automatically fall back to form data
			// unless we explicitly check, like in the PHP code.
			// However, nowadays webhooks are almost always JSON.
		}

		// Signature verification (Placeholder logic based on PHP comments)
		verified := tspaySvc.VerifyWebhookSignature(true, body, signature, timestamp)
		log.Printf("Incoming ctl @ deposit_callback verifiedSignature = %v", verified)

		// if !verified {
		// 	utils.WriteJSON(w, http.StatusUnauthorized, utils.Response{
		// 		Status:  "error",
		// 		Message: "Invalid signature",
		// 	})
		// 	return
		// }

		url := "https://api.allpayhub.com/incoming/tspay_deposit_callback"
		log.Printf("incoming Ctrl @ deposit_callback url = %s", url)

		// Forward to merchant service
		go incomingSvc.CallMerchant(url, data)

		utils.WriteJSON(w, http.StatusOK, utils.Response{
			Status:  "success",
			Message: "Deposit Callback received",
		})
	}
}
