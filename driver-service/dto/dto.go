package dto

import "bitaksi_taxihub/model"

// POST
type CreateDriverRequest struct {
	FirstName string  `json:"firstName" validate:"required"`
	LastName  string  `json:"lastName" validate:"required"`
	Plate     string  `json:"plate" validate:"required"`
	TaxiType  string  `json:"taxiType" validate:"required,oneof=sari buyuk luks"`
	CarBrand  string  `json:"carBrand" validate:"required"`
	CarModel  string  `json:"carModel" validate:"required"`
	Lat       float64 `json:"lat" validate:"required,gte=-90,lte=90"`
	Lon       float64 `json:"lon" validate:"required,gte=-180,lte=180"`
}

// PUT
type UpdateDriverRequest struct {
	FirstName *string  `json:"firstName" validate:"omitempty"`
	LastName  *string  `json:"lastName" validate:"omitempty"`
	Plate     *string  `json:"plate" validate:"omitempty"`
	TaxiType  *string  `json:"taxiType" validate:"omitempty,oneof=sari buyuk luks"`
	CarBrand  *string  `json:"carBrand" validate:"omitempty"`
	CarModel  *string  `json:"carModel" validate:"omitempty"`
	Lat       *float64 `json:"lat" validate:"omitempty,gte=-90,lte=90"`
	Lon       *float64 `json:"lon" validate:"omitempty,gte=-180,lte=180"`
}

// GET
type NearbyDriverResponse struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	Plate      string  `json:"plate"`
	DistanceKm float64 `json:"distanceKm"`
}

// GET
type ListDrivers struct {
	Items    []model.Driver `json:"items"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Total    int64          `json:"total"`
}
