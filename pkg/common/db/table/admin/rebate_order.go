package admin

import (
	"context"
	"time"

	"github.com/openimsdk/tools/db/pagination"
)

type RebateOrder struct {
	OrderID     string    `bson:"order_id"` // business ID (alphanumeric)
	Status      int32     `bson:"status"`   // 0 unread, 1 read, 2 approved, 3 rejected, 4 completed
	ReadAt      time.Time `bson:"read_at,omitempty"`
	ProcessedAt time.Time `bson:"processed_at,omitempty"`

	ProcurementOrderID  string    `bson:"procurement_order_id,omitempty"` // optional source procurement order
	ContractNumber      string    `bson:"contract_number"`                // contract no.
	InvoiceAmount       int64     `bson:"invoice_amount"`                 // cents
	InvoiceDate         time.Time `bson:"invoice_date"`                   // invoice date
	CustomerName        string    `bson:"customer_name"`                  // customer legal name
	CustomerBankAccount string    `bson:"customer_bank_account"`          // masked in UI as needed
	Remark              string    `bson:"remark,omitempty"`

	CreateTime time.Time `bson:"create_time"`
}

func (RebateOrder) TableName() string { return "rebate_order" }

type RebateOrderInterface interface {
	Add(ctx context.Context, ro *RebateOrder) error
	GetByOrderID(ctx context.Context, orderID string) (*RebateOrder, error)
	Update(ctx context.Context, orderID string, update map[string]any) error
	Delete(ctx context.Context, orderID string) error

	// Search by keyword (orderID / contractNumber / customerName), optional status and time range
	Search(ctx context.Context, keyword string, status *int32, start, end *time.Time, pg pagination.Pagination) (int64, []*RebateOrder, error)

	// Count by status (e.g., unread=0)
	CountByStatus(ctx context.Context, status int32) (int64, error)
}
