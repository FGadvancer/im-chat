package admin

import (
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	adminpb "github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/chat/pkg/util"
)

func ConvertEnterpriseToPB(enterprise *admindb.EnterpriseInfo) *adminpb.EnterpriseInfo {
	return &adminpb.EnterpriseInfo{
		EnterpriseID:          enterprise.EnterpriseID.Hex(),
		Name:                  enterprise.Name,
		Logo:                  enterprise.Logo,
		Website:               enterprise.Website,
		IsEligibleForCashback: enterprise.IsEligibleForCashback,
		Tags:                  enterprise.Tags,
		TagsTypes:             enterprise.TagsTypes,
		Address:               enterprise.Address,
		PhoneNumber:           enterprise.PhoneNumber,
		Email:                 enterprise.Email,
		Invoice:               enterprise.Invoice,
		Remark:                enterprise.Remark,
		CreateTime:            enterprise.CreateTime.UnixMilli(),
		Contacts:              enterprise.Contacts,
	}
}

func convertProcurementToPB(in *admindb.ProcurementOrder) *adminpb.ProcurementOder {
	if in == nil {
		return nil
	}
	return &adminpb.ProcurementOder{
		OrderID:         in.OrderID,
		Status:          in.Status,
		ReadAt:          util.ToMillis(in.ReadAt),
		ProcessedAt:     util.ToMillis(in.ProcessedAt),
		Vendor:          in.Vendor,
		Integrator:      in.Integrator,
		EndCustomer:     in.EndCustomer,
		PurchaseAmount:  in.PurchaseAmount,
		PaymentTermDays: in.PaymentTermDays,
		OriginOrgName:   in.OriginOrgName,
		ContactName:     in.ContactName,
		ContactPhone:    in.ContactPhone,
		ContactEmail:    in.ContactEmail,
		Remark:          in.Remark,
		CreateTime:      util.ToMillis(in.CreateTime),
	}
}

func convertRebateToPB(in *admindb.RebateOrder) *adminpb.RebateOder {
	if in == nil {
		return nil
	}
	return &adminpb.RebateOder{
		OrderID:             in.OrderID,
		Status:              in.Status,
		ReadAt:              util.ToMillis(in.ReadAt),
		ProcessedAt:         util.ToMillis(in.ProcessedAt),
		ProcurementOrderID:  in.ProcurementOrderID,
		ContractNumber:      in.ContractNumber,
		InvoiceAmount:       in.InvoiceAmount,
		InvoiceDate:         util.ToMillis(in.InvoiceDate),
		CustomerName:        in.CustomerName,
		CustomerBankAccount: in.CustomerBankAccount,
		Remark:              in.Remark,
		CreateTime:          util.ToMillis(in.CreateTime),
	}
}
