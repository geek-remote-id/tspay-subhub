package models

import "google.golang.org/protobuf/types/known/timestamppb"

// Category represents a category in the cashier system
type Category struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	DeletedAt   *timestamppb.Timestamp `json:"deleted_at"`
}
