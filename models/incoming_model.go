package models

// DepositCallback represents the data sent by the tspay webhook
type DepositCallback struct {
	Type          string      `json:"type"`
	TransactionID string      `json:"transaction_id"`
	Reference     string      `json:"reference"`
	Amount        float64     `json:"amount"`
	Fee           float64     `json:"fee"`
	NetAmount     float64     `json:"net_amount"`
	Currency      string      `json:"currency"`
	Chain         interface{} `json:"chain"`
	Status        string      `json:"status"`
	WalletAddress interface{} `json:"wallet_address"`
	Kind          string      `json:"kind"`
	Method        string      `json:"method"`
}
