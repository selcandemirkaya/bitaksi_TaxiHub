package repository

import (
	"bitaksi_taxihub/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//db operasyonları

type DriverRepository interface {
	//insert
	Insert(ctx context.Context, driver *model.Driver) (primitive.ObjectID, error)
	//update
	UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) error
	//List
	List(ctx context.Context, skip, limit int64) ([]model.Driver, error)
	//count
	Count(ctx context.Context) (int64, error)
	//findbytaxi
	FindByTaxiType(ctx context.Context, taxiType string) ([]model.Driver, error)
}

type driverRepositoryMongo struct {
	col *mongo.Collection
}

func NewDriverRepositoryMongo(col *mongo.Collection) DriverRepository {
	return &driverRepositoryMongo{col: col}
}

func (r *driverRepositoryMongo) Insert(ctx context.Context, d *model.Driver) (primitive.ObjectID, error) {
	res, err := r.col.InsertOne(ctx, d)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, mongo.ErrNilDocument
	}
	return id, nil
}

// UPDATE (partial)
func (r *driverRepositoryMongo) UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	filter := bson.M{"_id": id}
	updateDoc := bson.M{"$set": update}

	_, err := r.col.UpdateOne(ctx, filter, updateDoc)
	return err
}

// LIST (pagination)
func (r *driverRepositoryMongo) List(ctx context.Context, skip, limit int64) ([]model.Driver, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit)

	cur, err := r.col.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var drivers []model.Driver
	if err := cur.All(ctx, &drivers); err != nil {
		return nil, err
	}

	return drivers, nil
}

// COUNT (total)
func (r *driverRepositoryMongo) Count(ctx context.Context) (int64, error) {
	total, err := r.col.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	return total, nil
}

// FIND BY TAXI TYPE
func (r *driverRepositoryMongo) FindByTaxiType(ctx context.Context, taxiType string) ([]model.Driver, error) {
	filter := bson.M{"taxiType": taxiType}

	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var drivers []model.Driver
	if err := cur.All(ctx, &drivers); err != nil {
		return nil, err
	}

	return drivers, nil
}
