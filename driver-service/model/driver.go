package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// bson --> binary json(mongo)

type Driver struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"` //boşsa mongo üretsin
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Plate     string             `json:"plate" bson:"plate"`
	TaxiType  string             `json:"taxiType" bson:"taxiType"`
	CarBrand  string             `json:"carBrand" bson:"carBrand"`
	CarModel  string             `json:"carModel" bson:"carModel"`
	Lat       float64            `json:"lat" bson:"lat"` //latitude
	Lon       float64            `json:"lon" bson:"lon"` //longitude
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
