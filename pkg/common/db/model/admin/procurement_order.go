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

func NewProcurementOrder(db *mongo.Database) (admindb.ProcurementOrderInterface, error) {
	coll := db.Collection((admindb.ProcurementOrder{}).TableName())
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "order_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}, {Key: "create_time", Value: -1}},
		},
		{
			Keys:    bson.D{{Key: "create_time", Value: -1}},
			Options: options.Index().SetPartialFilterExpression(bson.M{"status": 0}).SetName("idx_unread_create_time"),
		},
		{
			Keys: bson.D{{Key: "contact_phone", Value: 1}},
		},
	})
	if err != nil {
		return nil, err
	}
	return &ProcurementOrderMgo{coll: coll}, nil
}

type ProcurementOrderMgo struct {
	coll *mongo.Collection
}

func (p *ProcurementOrderMgo) Add(ctx context.Context, po *admindb.ProcurementOrder) error {
	if po.CreateTime.IsZero() {
		po.CreateTime = time.Now()
	}
	return mongoutil.InsertMany(ctx, p.coll, []*admindb.ProcurementOrder{po})
}

func (p *ProcurementOrderMgo) GetByOrderID(ctx context.Context, orderID string) (*admindb.ProcurementOrder, error) {
	return mongoutil.FindOne[*admindb.ProcurementOrder](ctx, p.coll, bson.M{"order_id": orderID})
}

func (p *ProcurementOrderMgo) Update(ctx context.Context, orderID string, update map[string]any) error {
	if len(update) == 0 {
		return nil
	}

	// Build $set document from input update map
	setDoc := bson.M{}
	for k, v := range update {
		setDoc[k] = v
	}
	if raw, ok := update["status"]; ok {
		now := time.Now()
		switch v := raw.(type) {
		case int32:
			switch v {
			case 1: // read
				setDoc["read_at"] = now
			case 2, 3, 4: // processed-type states
				setDoc["processed_at"] = now
			}

		}
	}
	return mongoutil.UpdateOne(ctx, p.coll, bson.M{"order_id": orderID}, bson.M{"$set": setDoc}, true)
}

func (p *ProcurementOrderMgo) Delete(ctx context.Context, orderID string) error {
	return mongoutil.DeleteMany(ctx, p.coll, bson.M{"order_id": orderID})
}

func (p *ProcurementOrderMgo) Search(
	ctx context.Context,
	keyword string,
	status *int32,
	start, end *time.Time,
	pg pagination.Pagination,
) (int64, []*admindb.ProcurementOrder, error) {

	filter := bson.M{}
	if status != nil {
		filter["status"] = *status
	}
	if start != nil || end != nil {
		timeCond := bson.M{}
		if start != nil {
			timeCond["$gte"] = *start
		}
		if end != nil {
			timeCond["$lte"] = *end
		}
		filter["create_time"] = timeCond
	}
	if kw := regexp.QuoteMeta(keyword); kw != "" {
		filter["$or"] = []bson.M{
			{"order_id": bson.M{"$regex": kw, "$options": "i"}},
			{"contact_phone": bson.M{"$regex": kw, "$options": "i"}},
		}
	}

	return mongoutil.FindPage[*admindb.ProcurementOrder](
		ctx,
		p.coll,
		filter,
		pg,
		options.Find().SetSort(p.sort()),
	)
}

func (p *ProcurementOrderMgo) UpdateStatus(
	ctx context.Context,
	orderID string,
	to int32,
	expectedFrom *int32,
	set map[string]any,
	unset []string,
) error {
	filter := bson.M{"order_id": orderID}
	if expectedFrom != nil {
		filter["status"] = *expectedFrom
	}

	if set == nil {
		set = make(map[string]any, 1)
	}
	set["status"] = to

	upd := bson.M{"$set": set}
	if len(unset) > 0 {
		unsetMap := bson.M{}
		for _, k := range unset {
			unsetMap[k] = ""
		}
		upd["$unset"] = unsetMap
	}

	return mongoutil.UpdateOne(ctx, p.coll, filter, upd, false)
}
func (p *ProcurementOrderMgo) CountByStatus(ctx context.Context, status int32) (int64, error) {
	return mongoutil.Count(ctx, p.coll, bson.M{"status": status})
}

func (p *ProcurementOrderMgo) sort() any {
	return bson.D{{Key: "create_time", Value: -1}}
}
