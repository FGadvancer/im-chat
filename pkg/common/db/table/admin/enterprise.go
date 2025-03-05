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

	"github.com/openimsdk/tools/db/pagination"
)

type EnterpriseInfo struct {
	EnterpriseID          primitive.ObjectID `bson:"_id"`
	Name                  string             `bson:"name"`
	Logo                  string             `bson:"logo"`
	Website               string             `bson:"website"`
	IsEligibleForCashback bool               `bson:"is_eligible_for_cashback"`
	Tags                  []string           `bson:"tags"`
	Address               string             `bson:"address"`
	PhoneNumber           string             `bson:"phone_number"`
	Email                 string             `bson:"email"`
	CreateTime            time.Time          `bson:"create_time"`
}

func (EnterpriseInfo) TableName() string {
	return "enterprise_info"
}

type EnterpriseInfoInterface interface {
	Add(ctx context.Context, info *EnterpriseInfo) error
	Get(ctx context.Context, enterpriseID primitive.ObjectID) (*EnterpriseInfo, error)
	Update(ctx context.Context, id primitive.ObjectID, update map[string]any) error
	Delete(ctx context.Context, enterpriseID primitive.ObjectID) error
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*EnterpriseInfo, error)
}
