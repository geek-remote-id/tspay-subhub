package models

import "google.golang.org/protobuf/types/known/timestamppb"

// Category represents a category in the cashier system
type Category struct {
	ID          int                    `json:"id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	DeletedAt   *timestamppb.Timestamp `json:"deleted_at,omitempty"`
}
