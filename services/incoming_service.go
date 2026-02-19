package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type IncomingService struct{}

func NewIncomingService() *IncomingService {
	return &IncomingService{}
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
