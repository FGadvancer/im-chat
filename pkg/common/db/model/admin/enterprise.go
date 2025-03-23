package admin

import (
	"context"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
)

func NewEnterpriseInfo(db *mongo.Database) (admindb.EnterpriseInfoInterface, error) {
	coll := db.Collection("enterprise_info")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "create_time", Value: -1},
			},
		},
		{Keys: bson.D{
			{Key: "tags", Value: 1},
		},
		},
	})
	if err != nil {
		return nil, err
	}
	return &EnterpriseInfoMgo{coll: coll}, nil
}

type EnterpriseInfoMgo struct {
	coll *mongo.Collection
}

func (a *EnterpriseInfoMgo) Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*admindb.EnterpriseInfo, error) {
	escapedKeyword := regexp.QuoteMeta(keyword)
	filter := bson.M{}
	if escapedKeyword != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": escapedKeyword, "$options": "i"}},
				{"tags": bson.M{"$regex": escapedKeyword, "$options": "i"}},
			},
		}
	}
	return mongoutil.FindPage[*admin.EnterpriseInfo](ctx, a.coll, filter, pagination, options.Find().SetSort(a.sort()))
}

func (a *EnterpriseInfoMgo) Add(ctx context.Context, info *admindb.EnterpriseInfo) error {
	return mongoutil.InsertMany(ctx, a.coll, []*admin.EnterpriseInfo{info})
}

func (a *EnterpriseInfoMgo) Get(ctx context.Context, enterpriseID primitive.ObjectID) (*admindb.EnterpriseInfo, error) {
	return mongoutil.FindOne[*admin.EnterpriseInfo](ctx, a.coll, bson.M{"_id": enterpriseID})
}

func (a *EnterpriseInfoMgo) Update(ctx context.Context, id primitive.ObjectID, update map[string]any) error {
	if len(update) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, a.coll, bson.M{"_id": id}, bson.M{"$set": update}, true)
}

func (a *EnterpriseInfoMgo) Delete(ctx context.Context, enterpriseID primitive.ObjectID) error {
	return mongoutil.DeleteMany(ctx, a.coll, bson.M{"_id": enterpriseID})
}

func (a *EnterpriseInfoMgo) sort() any {
	return bson.D{{"create_time", -1}}
}
