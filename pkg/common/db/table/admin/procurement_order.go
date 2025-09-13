package admin

import (
	"context"
	"time"

	"github.com/openimsdk/tools/db/pagination"
)

type ProcurementOrder struct {
	OrderID     string    `bson:"order_id"`               // Business order identifier (alphanumeric, unique)
	Status      int32     `bson:"status"`                 // Status: 0=unread, 1=read, 2=approved, 3=rejected, 4=completed (extensible)
	ReadAt      time.Time `bson:"read_at,omitempty"`      // Timestamp when marked as read
	ProcessedAt time.Time `bson:"processed_at,omitempty"` // Timestamp when processing/decision occurred

	// Business fields
	Vendor          string    `bson:"vendor"`                  // Manufacturer
	Integrator      string    `bson:"integrator"`              // Integrator/Reseller
	EndCustomer     string    `bson:"end_customer"`            // End customer (purchasing party)
	PurchaseAmount  int64     `bson:"purchase_amount"`         // Purchase amount in cents (int64 to avoid floating-point errors)
	PaymentTermDays int32     `bson:"payment_term_days"`       // Payment term in days
	OriginOrgName   string    `bson:"origin_org_name"`         // Originating organization name
	ContactName     string    `bson:"contact_name"`            // Contact person name
	ContactPhone    string    `bson:"contact_phone"`           // Contact phone number
	ContactEmail    string    `bson:"contact_email,omitempty"` // Contact email (optional)
	Remark          string    `bson:"remark,omitempty"`        // Additional notes (optional)
	CreateTime      time.Time `bson:"create_time"`             // Creation timestamp
}

func (ProcurementOrder) TableName() string {
	return "procurement_order"
}

type ProcurementOrderInterface interface {
	Add(ctx context.Context, po *ProcurementOrder) error
	GetByOrderID(ctx context.Context, orderID string) (*ProcurementOrder, error)
	Update(ctx context.Context, orderID string, update map[string]any) error
	Delete(ctx context.Context, orderID string) error

	// Search orders by keyword (matches orderID or contactPhone), with optional status
	// and creation-time range filters. Results are paginated.
	Search(ctx context.Context, keyword string, status *int32, start, end *time.Time, pg pagination.Pagination) (int64, []*ProcurementOrder, error)

	UpdateStatus(ctx context.Context, orderID string, to int32, expectedFrom *int32, set map[string]any, unset []string) error
	CountByStatus(ctx context.Context, status int32) (int64, error)
}
