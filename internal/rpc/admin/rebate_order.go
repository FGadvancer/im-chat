package admin

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	adminpb "github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

// Add

func (o *adminServer) AddRebateOder(ctx context.Context, req *adminpb.AddRebateOderReq) (*adminpb.AddRebateOderResp, error) {

	orderID, err := o.newReadableRebateOrderID(ctx)
	if err != nil {
		log.ZWarn(ctx, "generate rebate order id failed", err)
		return nil, err
	}
	val := &admindb.RebateOrder{
		OrderID:             orderID,
		Status:              0,
		ProcurementOrderID:  req.ProcurementOrderID,
		ContractNumber:      req.ContractNumber,
		InvoiceAmount:       req.InvoiceAmount,
		InvoiceDate:         time.UnixMilli(req.InvoiceDate),
		CustomerName:        req.CustomerName,
		CustomerBankAccount: req.CustomerBankAccount,
		Remark:              req.Remark,
		CreateTime:          time.Now(),
	}
	if err := o.Database.AddRebateOrder(ctx, val); err != nil {
		log.ZWarn(ctx, "add rebate order failed", err)
		return nil, err
	}
	return &adminpb.AddRebateOderResp{OrderID: orderID}, nil
}

// Update

func (o *adminServer) UpdateRebateOder(ctx context.Context, req *adminpb.UpdateRebateOderReq) (*adminpb.UpdateRebateOderResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if req.OrderID == "" {
		return nil, errs.ErrArgs.WrapMsg("orderID is empty")
	}
	update := make(map[string]any)

	if req.ProcurementOrderID != nil {
		update["procurement_order_id"] = *req.ProcurementOrderID
	}
	if req.ContractNumber != nil {
		update["contract_number"] = *req.ContractNumber
	}
	if req.InvoiceAmount != nil {
		update["invoice_amount"] = *req.InvoiceAmount
	}
	if req.InvoiceDate != nil {
		update["invoice_date"] = time.UnixMilli(*req.InvoiceDate)
	}
	if req.CustomerName != nil {
		update["customer_name"] = *req.CustomerName
	}
	if req.CustomerBankAccount != nil {
		update["customer_bank_account"] = *req.CustomerBankAccount
	}
	if req.Remark != nil {
		update["remark"] = *req.Remark
	}
	if req.Status != nil {
		update["status"] = *req.Status
	}

	if len(update) == 0 {
		return &adminpb.UpdateRebateOderResp{}, nil
	}
	if err := o.Database.UpdateRebateOrder(ctx, req.OrderID, update); err != nil {
		return nil, err
	}
	return &adminpb.UpdateRebateOderResp{}, nil
}

// Delete

func (o *adminServer) DeleteRebateOder(ctx context.Context, req *adminpb.DeleteRebateOderReq) (*adminpb.DeleteRebateOderResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if req.OrderID == "" {
		return nil, errs.ErrArgs.WrapMsg("orderID is empty")
	}
	if err := o.Database.DeleteRebateOrder(ctx, req.OrderID); err != nil {
		return nil, err
	}
	return &adminpb.DeleteRebateOderResp{}, nil
}

// Get

func (o *adminServer) GetRebateOder(ctx context.Context, req *adminpb.GetRebateOderReq) (*adminpb.GetRebateOderResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if req.OrderID == "" {
		return nil, errs.ErrArgs.WrapMsg("orderID is empty")
	}
	ro, err := o.Database.GetRebateOrder(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	return &adminpb.GetRebateOderResp{
		RebateOrder: convertRebateToPB(ro),
	}, nil
}

// Query

func (o *adminServer) QueryRebateOderList(ctx context.Context, req *adminpb.QueryRebateOderListReq) (*adminpb.QueryRebateOderListResp, error) {
	var (
		statusPtr *int32
		startPtr  *time.Time
		endPtr    *time.Time
	)
	if req.Status != nil {
		status := *req.Status
		statusPtr = &status
	}
	if req.StartCreateTime != nil && *req.StartCreateTime > 0 {
		t := time.UnixMilli(*req.StartCreateTime)
		startPtr = &t
	}
	if req.EndCreateTime != nil && *req.EndCreateTime > 0 {
		t := time.UnixMilli(*req.EndCreateTime)
		endPtr = &t
	}

	total, items, err := o.Database.SearchRebateOrder(ctx, req.Keyword, statusPtr, startPtr, endPtr, req.Pagination)
	if err != nil {
		return nil, err
	}

	resp := &adminpb.QueryRebateOderListResp{
		Total:        int32(total),
		RebateOrders: make([]*adminpb.RebateOder, 0, len(items)),
	}
	for _, it := range items {
		resp.RebateOrders = append(resp.RebateOrders, convertRebateToPB(it))
	}
	return resp, nil
}

// Count

func (o *adminServer) CountRebateOderByStatus(ctx context.Context, req *adminpb.CountRebateOderByStatusReq) (*adminpb.CountRebateOderByStatusResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	cnt, err := o.Database.CountRebateOrderByStatus(ctx, req.Status)
	if err != nil {
		return nil, err
	}
	return &adminpb.CountRebateOderByStatusResp{Count: int32(cnt)}, nil
}

// Helpers

// Replace with your atomic daily sequence generator if desired.
func (o *adminServer) newReadableRebateOrderID(_ context.Context) (string, error) {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	day := time.Now().Format("20060102")
	n := 6
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		v, err := rand.Int(rand.Reader, big.NewInt(36))
		if err != nil {
			return "", err
		}
		buf[i] = alphabet[v.Int64()]
	}
	return fmt.Sprintf("RB-%s-%s", day, string(buf)), nil
}
