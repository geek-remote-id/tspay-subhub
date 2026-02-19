package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/geek-remote-id/tspay-subhub/utils"
	"github.com/spf13/viper"
)

type TspayService struct {
	toleranceSeconds     int64
	webhookSecretDeposit string
	webhookSecretPayout  string
}

func NewTspayService() *TspayService {
	tolerance := viper.GetInt64("TSPAY_WEBHOOK_TOLERANCE")
	if tolerance == 0 {
		tolerance = 300 // default 5 minutes
	}

	return &TspayService{
		toleranceSeconds:     tolerance,
		webhookSecretDeposit: viper.GetString("TSPAY_WEBHOOK_SECRET_DEPOSIT"),
		webhookSecretPayout:  viper.GetString("TSPAY_WEBHOOK_SECRET_PAYOUT"),
	}
}

// VerifyWebhookSignature verifies the signature of the incoming webhook
func (s *TspayService) VerifyWebhookSignature(isDeposit bool, payload []byte, signature string, timestamp string) bool {

	currentTime := time.Now().Unix()
	webhookTime, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		utils.LogEvent("events", "Warning: Could not parse timestamp '%s' as number, skipping tolerance check.", timestamp)
	} else {
		// Only check tolerance if we have a valid numeric timestamp
		if math.Abs(float64(currentTime-webhookTime)) > float64(s.toleranceSeconds) {
			utils.LogEvent("errors", "Webhook timestamp outside tolerance window: current=%d, webhook=%d", currentTime, webhookTime)
			return false
		}
	}

	// Remove 'sha256=' prefix if present
	signature = strings.TrimPrefix(signature, "sha256=")

	// Determine which secret key to use
	secret := s.webhookSecretDeposit
	if !isDeposit {
		secret = s.webhookSecretPayout
	}

	// Compute expected signature
	message := fmt.Sprintf("%s.%s", timestamp, string(payload))
	utils.LogEvent("events", "Incoming ctl @ deposit_callback message = %s", message)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	utils.LogEvent("events", "Incoming ctl @ deposit_callback secret = %s", secret)
	utils.LogEvent("events", "Incoming ctl @ deposit_callback expectedSignature = %s", expectedSignature)
	utils.LogEvent("events", "Incoming ctl @ deposit_callback signature = %s", signature)

	// Constant-time comparison
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
