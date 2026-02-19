package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/geek-remote-id/tspay-subhub/models"
)

type IncomingService struct {
	tspaySvc *TspayService
}

func NewIncomingService(tspaySvc *TspayService) *IncomingService {
	return &IncomingService{
		tspaySvc: tspaySvc,
	}
}

// ProcessDepositCallback handles the core logic of verifying and forwarding a deposit callback
func (s *IncomingService) ProcessDepositCallback(body []byte, signature, timestamp string) error {
	// 1. Verify Signature
	if !s.tspaySvc.VerifyWebhookSignature(true, body, signature, timestamp) {
		return fmt.Errorf("invalid signature")
	}

	// 2. Unmarshal data
	var data models.DepositCallback
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// 3. Log and Forward
	log.Printf("Processing deposit callback for Transaction ID: %s", data.TransactionID)
	url := "https://api.allpayhub.com/incoming/tspay_deposit_callback"

	// Forward asynchronously to avoid blocking the webhook response
	go s.CallMerchant(url, data)

	return nil
}

// CallMerchant forwards the callback data to the merchant's URL
func (s *IncomingService) CallMerchant(url string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	log.Printf("Forwarding callback to merchant: %s", url)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error calling merchant: %v", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("Merchant responded with status: %s", resp.Status)
	return nil
}
