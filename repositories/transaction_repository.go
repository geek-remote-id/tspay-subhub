package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"strings"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction creates a new transaction with its details
func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	// Step 1: Validate all products and check stock availability
	type productInfo struct {
		name  string
		price int
		stock int
	}
	productData := make(map[int]productInfo)

	for _, item := range items {
		var name string
		var price, stock int

		err := tx.QueryRow("SELECT name, price, stock FROM product WHERE id = $1 AND deleted_at IS NULL", item.ProductID).Scan(&name, &price, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Validate stock availability
		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product '%s' (available: %d, requested: %d)", name, stock, item.Quantity)
		}

		productData[item.ProductID] = productInfo{
			name:  name,
			price: price,
			stock: stock,
		}
	}

	// Step 2: Calculate total and prepare details
	for _, item := range items {
		product := productData[item.ProductID]
		subtotal := product.price * item.Quantity
		totalAmount += subtotal

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: product.name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Step 3: Update stock for all products
	for _, item := range items {
		_, err = tx.Exec("UPDATE product SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}
	}

	// Step 4: Insert transaction record
	var transactionID int
	var createdAt, deletedAt sql.NullTime
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at, deleted_at", totalAmount).Scan(&transactionID, &createdAt, &deletedAt)
	if err != nil {
		return nil, err
	}

	// Step 5: Batch insert transaction details
	if len(details) > 0 {
		valueStrings := make([]string, 0, len(details))
		valueArgs := make([]interface{}, 0, len(details)*4)

		for i, detail := range details {
			details[i].TransactionID = transactionID
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)",
				i*4+1, i*4+2, i*4+3, i*4+4))
			valueArgs = append(valueArgs, transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
		}

		query := fmt.Sprintf("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES %s",
			strings.Join(valueStrings, ","))

		_, err = tx.Exec(query, valueArgs...)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	transaction := &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}

	// Database connection already handles timezone conversion
	// Timestamps are returned in Asia/Jakarta timezone (UTC+7)
	if createdAt.Valid {
		transaction.CreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
	}

	if deletedAt.Valid {
		transaction.DeletedAt = deletedAt.Time.Format("2006-01-02 15:04:05")
	}

	return transaction, nil
}
