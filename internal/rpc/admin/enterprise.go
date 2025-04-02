// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package admin

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	adminpb "github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

func (o *adminServer) AddEnterpriseInfo(ctx context.Context, req *adminpb.AddEnterpriseInfoReq) (*adminpb.AddEnterpriseInfoResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	id := primitive.NewObjectID()
	val := &admindb.EnterpriseInfo{
		EnterpriseID:          id,
		Name:                  req.Name,
		Logo:                  req.Logo,
		Website:               req.Website,
		IsEligibleForCashback: req.IsEligibleForCashback,
		Tags:                  req.Tags,
		Address:               req.Address,
		PhoneNumber:           req.PhoneNumber,
		Email:                 req.Email,
		CreateTime:            time.Now(),
	}
	if err := o.Database.AddEnterpriseInfo(ctx, val); err != nil {
		log.ZWarn(ctx, "add enterprise info failed", err)
		return nil, err
	}
	return &adminpb.AddEnterpriseInfoResp{EnterpriseID: id.Hex()}, nil
}
func (o *adminServer) UpdateEnterpriseInfo(ctx context.Context, req *adminpb.UpdateEnterpriseInfoReq) (*adminpb.UpdateEnterpriseInfoResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	oid, err := primitive.ObjectIDFromHex(req.EnterpriseID)
	if err != nil {
		return nil, errs.ErrArgs.WrapMsg("invalid id " + err.Error())
	}
	update := make(map[string]any)
	if req.Name != nil {
		update["name"] = *req.Name
	}
	if req.Logo != nil {
		update["logo"] = *req.Logo
	}
	if req.Website != nil {
		update["website"] = *req.Website
	}
	if req.IsEligibleForCashback != nil {
		update["is_eligible_for_cashback"] = *req.IsEligibleForCashback
	}
	if len(req.Tags) > 0 {
		update["tags"] = req.Tags
	}
	if req.ClearTags {
		update["tags"] = []string{}
	}
	if req.Address != nil {
		update["address"] = *req.Address
	}
	if req.PhoneNumber != nil {
		update["phone_number"] = *req.PhoneNumber
	}
	if req.Email != nil {
		update["email"] = *req.Email
	}
	if err := o.Database.UpdateEnterpriseInfo(ctx, oid, update); err != nil {
		return nil, err
	}
	return &adminpb.UpdateEnterpriseInfoResp{}, nil

}
func (o *adminServer) DeleteEnterpriseInfo(ctx context.Context, req *adminpb.DeleteEnterpriseInfoReq) (*adminpb.DeleteEnterpriseInfoResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	oid, err := primitive.ObjectIDFromHex(req.EnterpriseID)
	if err != nil {
		return nil, errs.ErrArgs.WrapMsg("invalid id " + err.Error())
	}

	if err := o.Database.DeleteEnterpriseInfo(ctx, oid); err != nil {
		return nil, err
	}
	return &adminpb.DeleteEnterpriseInfoResp{}, nil
}
func (o *adminServer) GetEnterpriseInfo(ctx context.Context, req *adminpb.GetEnterpriseInfoReq) (*adminpb.GetEnterpriseInfoResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	oid, err := primitive.ObjectIDFromHex(req.EnterpriseID)
	if err != nil {
		return nil, errs.ErrArgs.WrapMsg("invalid id " + err.Error())
	}
	enterprise, err := o.Database.GetEnterpriseInfo(ctx, oid)
	if err != nil {
		return nil, err
	}
	return &adminpb.GetEnterpriseInfoResp{
		Enterprise: &adminpb.EnterpriseInfo{
			EnterpriseID:          enterprise.EnterpriseID.Hex(),
			Name:                  enterprise.Name,
			Logo:                  enterprise.Logo,
			Website:               enterprise.Website,
			IsEligibleForCashback: enterprise.IsEligibleForCashback,
			Tags:                  enterprise.Tags,
			Address:               enterprise.Address,
			PhoneNumber:           enterprise.PhoneNumber,
			Email:                 enterprise.Email,
			CreateTime:            enterprise.CreateTime.UnixMilli(),
		},
	}, nil
}

func (o *adminServer) QueryEnterpriseList(ctx context.Context, req *adminpb.QueryEnterpriseListReq) (*adminpb.QueryEnterpriseListResp, error) {
	total, enterprises, err := o.Database.SearchEnterpriseInfo(ctx, req.NameKeyword, req.Pagination)
	if err != nil {
		return nil, err
	}
	resp := &adminpb.QueryEnterpriseListResp{Total: int32(total), Enterprises: make([]*adminpb.EnterpriseInfo, 0, len(enterprises))}
	for _, enterprise := range enterprises {
		resp.Enterprises = append(resp.Enterprises, &adminpb.EnterpriseInfo{
			EnterpriseID:          enterprise.EnterpriseID.Hex(),
			Name:                  enterprise.Name,
			Logo:                  enterprise.Logo,
			Website:               enterprise.Website,
			IsEligibleForCashback: enterprise.IsEligibleForCashback,
			Tags:                  enterprise.Tags,
			Address:               enterprise.Address,
			PhoneNumber:           enterprise.PhoneNumber,
			Email:                 enterprise.Email,
			CreateTime:            enterprise.CreateTime.UnixMilli(),
		})
	}
	return resp, nil
}
