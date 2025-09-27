package admin

import (
	"context"
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	adminpb "github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/chat/pkg/util"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"time"
)

// ---------- Add ----------

func (o *adminServer) AddProcurementOder(ctx context.Context, req *adminpb.AddProcurementOderReq) (*adminpb.AddProcurementOderResp, error) {

	orderID, err := util.NewReadableOrderID(ctx)
	if err != nil {
		log.ZWarn(ctx, "generate order id failed", err)
		return nil, err
	}

	val := &admindb.ProcurementOrder{
		OrderID:         orderID,
		Status:          0, // default unread
		Vendor:          req.Vendor,
		Integrator:      req.Integrator,
		EndCustomer:     req.EndCustomer,
		PurchaseAmount:  req.PurchaseAmount,
		PaymentTermDays: req.PaymentTermDays,
		OriginOrgName:   req.OriginOrgName,
		ContactName:     req.ContactName,
		ContactPhone:    req.ContactPhone,
		ContactEmail:    req.ContactEmail,
		Remark:          req.Remark,
		CreateTime:      time.Now(),
	}

	if err := o.Database.AddProcurementOrder(ctx, val); err != nil {
		log.ZWarn(ctx, "add procurement order failed", err)
		return nil, err
	}
	return &adminpb.AddProcurementOderResp{OrderID: orderID}, nil
}

// ---------- Update ----------

func (o *adminServer) UpdateProcurementOder(ctx context.Context, req *adminpb.UpdateProcurementOderReq) (*adminpb.UpdateProcurementOderResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if req.OrderID == "" {
		return nil, errs.ErrArgs.WrapMsg("orderID is empty")
	}

	update := make(map[string]any)
	// business fields
	if req.Vendor != nil {
		update["vendor"] = *req.Vendor
	}
	if req.Integrator != nil {
		update["integrator"] = *req.Integrator
	}
	if req.EndCustomer != nil {
		update["end_customer"] = *req.EndCustomer
	}
	if req.PurchaseAmount != nil {
		update["purchase_amount"] = *req.PurchaseAmount
	}
	if req.PaymentTermDays != nil {
		update["payment_term_days"] = *req.PaymentTermDays
	}
	if req.OriginOrgName != nil {
		update["origin_org_name"] = *req.OriginOrgName
	}
	if req.ContactName != nil {
		update["contact_name"] = *req.ContactName
	}
	if req.ContactPhone != nil {
		update["contact_phone"] = *req.ContactPhone
	}
	if req.ContactEmail != nil {
		update["contact_email"] = *req.ContactEmail
	}
	if req.Remark != nil {
		update["remark"] = *req.Remark
	}
	// status transition (DB layer stays generic)
	if req.Status != nil {
		update["status"] = *req.Status
	}

	if len(update) == 0 {
		// no-op update
		return &adminpb.UpdateProcurementOderResp{}, nil
	}
	if err := o.Database.UpdateProcurementOrder(ctx, req.OrderID, update); err != nil {
		return nil, err
	}
	return &adminpb.UpdateProcurementOderResp{}, nil
}

// ---------- Delete ----------

func (o *adminServer) DeleteProcurementOder(ctx context.Context, req *adminpb.DeleteProcurementOderReq) (*adminpb.DeleteProcurementOderResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if req.OrderID == "" {
		return nil, errs.ErrArgs.WrapMsg("orderID is empty")
	}
	if err := o.Database.DeleteProcurementOrder(ctx, req.OrderID); err != nil {
		return nil, err
	}
	return &adminpb.DeleteProcurementOderResp{}, nil
}

// ---------- Get ----------

func (o *adminServer) GetProcurementOder(ctx context.Context, req *adminpb.GetProcurementOderReq) (*adminpb.GetProcurementOderResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if req.OrderID == "" {
		return nil, errs.ErrArgs.WrapMsg("orderID is empty")
	}
	po, err := o.Database.GetProcurementOrder(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	return &adminpb.GetProcurementOderResp{
		ProcurementOrder: convertProcurementToPB(po),
	}, nil
}

// ---------- Query list ----------

func (o *adminServer) QueryProcurementOderList(ctx context.Context, req *adminpb.QueryProcurementOderListReq) (*adminpb.QueryProcurementOderListResp, error) {
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

	total, items, err := o.Database.SearchProcurementOrder(ctx, req.Keyword, statusPtr, startPtr, endPtr, req.Pagination)
	if err != nil {
		return nil, err
	}

	resp := &adminpb.QueryProcurementOderListResp{
		Total:             int32(total),
		ProcurementOrders: make([]*adminpb.ProcurementOder, 0, len(items)),
	}
	for _, it := range items {
		resp.ProcurementOrders = append(resp.ProcurementOrders, convertProcurementToPB(it))
	}
	return resp, nil
}

// ---------- Count by status ----------

func (o *adminServer) CountProcurementOderByStatus(ctx context.Context, req *adminpb.CountProcurementOderByStatusReq) (*adminpb.CountProcurementOderByStatusResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	cnt, err := o.Database.CountProcurementOrderByStatus(ctx, req.Status)
	if err != nil {
		return nil, err
	}
	return &adminpb.CountProcurementOderByStatusResp{Count: int32(cnt)}, nil
}
