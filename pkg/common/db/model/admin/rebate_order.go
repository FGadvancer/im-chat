package admin

import (
	"context"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
)

type RebateOrderMgo struct {
	coll *mongo.Collection
}

func NewRebateOrder(db *mongo.Database) (admindb.RebateOrderInterface, error) {
	coll := db.Collection((admindb.RebateOrder{}).TableName())
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{ // unique business ID
			Keys:    bson.D{{Key: "order_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{ // status + time for list pages
			Keys: bson.D{{Key: "status", Value: 1}, {Key: "create_time", Value: -1}},
		},
		{ // fast unread list/count
			Keys:    bson.D{{Key: "create_time", Value: -1}},
			Options: options.Index().SetPartialFilterExpression(bson.M{"status": 0}).SetName("idx_unread_create_time"),
		},
		{ // typical search fields
			Keys: bson.D{{Key: "contract_number", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "customer_name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "procurement_order_id", Value: 1}},
		},
	})
	if err != nil {
		return nil, err
	}
	return &RebateOrderMgo{coll: coll}, nil
}

func (r *RebateOrderMgo) Add(ctx context.Context, ro *admindb.RebateOrder) error {
	if ro.CreateTime.IsZero() {
		ro.CreateTime = time.Now()
	}
	return mongoutil.InsertMany(ctx, r.coll, []*admindb.RebateOrder{ro})
}

func (r *RebateOrderMgo) GetByOrderID(ctx context.Context, orderID string) (*admindb.RebateOrder, error) {
	return mongoutil.FindOne[*admindb.RebateOrder](ctx, r.coll, bson.M{"order_id": orderID})
}

func (r *RebateOrderMgo) Update(ctx context.Context, orderID string, update map[string]any) error {
	if len(update) == 0 {
		return nil
	}
	setDoc := bson.M{}
	for k, v := range update {
		setDoc[k] = v
	}
	// status-driven timestamps
	if raw, ok := update["status"]; ok {
		now := time.Now()
		switch v := raw.(type) {
		case int32:
			switch v {
			case 1: // read
				setDoc["read_at"] = now
			case 2, 3, 4: // processed-type states
				setDoc["processed_at"] = now
			case 0: // unread
				return mongoutil.UpdateOne(ctx, r.coll,
					bson.M{"order_id": orderID},
					bson.M{"$set": setDoc, "$unset": bson.M{"read_at": "", "processed_at": ""}},
					true,
				)
			}
		}
	}
	return mongoutil.UpdateOne(ctx, r.coll, bson.M{"order_id": orderID}, bson.M{"$set": setDoc}, true)
}

func (r *RebateOrderMgo) Delete(ctx context.Context, orderID string) error {
	return mongoutil.DeleteMany(ctx, r.coll, bson.M{"order_id": orderID})
}

func (r *RebateOrderMgo) Search(
	ctx context.Context,
	keyword string,
	status *int32,
	start, end *time.Time,
	pg pagination.Pagination,
) (int64, []*admindb.RebateOrder, error) {

	filter := bson.M{}
	if status != nil {
		filter["status"] = *status
	}
	if start != nil || end != nil {
		tm := bson.M{}
		if start != nil {
			tm["$gte"] = *start
		}
		if end != nil {
			tm["$lte"] = *end
		}
		filter["create_time"] = tm
	}
	if kw := regexp.QuoteMeta(keyword); kw != "" {
		filter["$or"] = []bson.M{
			{"order_id": bson.M{"$regex": kw, "$options": "i"}},
			{"contract_number": bson.M{"$regex": kw, "$options": "i"}},
			{"customer_name": bson.M{"$regex": kw, "$options": "i"}},
		}
	}

	return mongoutil.FindPage[*admindb.RebateOrder](
		ctx,
		r.coll,
		filter,
		pg,
		options.Find().SetSort(bson.D{{Key: "create_time", Value: -1}}),
	)
}

func (r *RebateOrderMgo) CountByStatus(ctx context.Context, status int32) (int64, error) {
	return mongoutil.Count(ctx, r.coll, bson.M{"status": status})
}
