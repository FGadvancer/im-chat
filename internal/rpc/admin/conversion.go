package admin

import (
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	adminpb "github.com/openimsdk/chat/pkg/protocol/admin"
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
