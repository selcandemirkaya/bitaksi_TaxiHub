package service

import (
	"context"
	"time"

	"bitaksi_taxihub/dto"
	"bitaksi_taxihub/model"
	"bitaksi_taxihub/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"sort"

	"github.com/umahmood/haversine"
)

type DriverService interface {
	CreateDriver(ctx context.Context, req dto.CreateDriverRequest) (primitive.ObjectID, error)
	UpdateDriver(ctx context.Context, id primitive.ObjectID, req dto.UpdateDriverRequest) error
	ListDrivers(ctx context.Context, page, pageSize int64) (dto.ListDrivers, error)
	NearbyDrivers(ctx context.Context, lat, lon float64, taxiType string) ([]dto.NearbyDriverResponse, error)
}

type driverService struct {
	repo repository.DriverRepository
}

func NewDriverService(repo repository.DriverRepository) DriverService {
	return &driverService{repo: repo}
}

func (s *driverService) CreateDriver(ctx context.Context, req dto.CreateDriverRequest) (primitive.ObjectID, error) {
	now := time.Now()

	driver := &model.Driver{
		ID:        primitive.NewObjectID(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Plate:     req.Plate,
		TaxiType:  req.TaxiType,
		CarBrand:  req.CarBrand,
		CarModel:  req.CarModel,
		Lat:       req.Lat,
		Lon:       req.Lon,
		CreatedAt: now,
		UpdatedAt: now,
	}

	id, err := s.repo.Insert(ctx, driver)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id, nil
}

func (s *driverService) UpdateDriver(ctx context.Context, id primitive.ObjectID, req dto.UpdateDriverRequest) error {
	update := bson.M{}
	if req.FirstName != nil {
		update["firstName"] = *req.FirstName
	}
	if req.LastName != nil {
		update["lastName"] = *req.LastName
	}
	if req.Plate != nil {
		update["plate"] = *req.Plate
	}
	if req.TaxiType != nil {
		update["taxiType"] = *req.TaxiType
	}
	if req.CarBrand != nil {
		update["carBrand"] = *req.CarBrand
	}
	if req.CarModel != nil {
		update["carModel"] = *req.CarModel
	}
	if req.Lat != nil {
		update["lat"] = *req.Lat
	}
	if req.Lon != nil {
		update["lon"] = *req.Lon
	}

	if len(update) == 0 {
		return nil
	}

	update["updatedAt"] = time.Now()
	return s.repo.UpdateByID(ctx, id, update)

}

func (s *driverService) ListDrivers(ctx context.Context, page, pageSize int64) (dto.ListDrivers, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	skip := (page - 1) * pageSize
	drivers, err := s.repo.List(ctx, skip, pageSize)
	if err != nil {
		return dto.ListDrivers{}, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return dto.ListDrivers{}, err
	}
	return dto.ListDrivers{
		Items:    drivers,
		Page:     int(page),
		PageSize: int(pageSize),
		Total:    total,
	}, nil
}

func (s *driverService) NearbyDrivers(ctx context.Context, lat, lon float64, taxiType string) ([]dto.NearbyDriverResponse, error) {
	drivers, err := s.repo.FindByTaxiType(ctx, taxiType)
	if err != nil {
		return nil, err
	}
	userCoord := haversine.Coord{Lat: lat, Lon: lon}

	var nearby []dto.NearbyDriverResponse

	for _, d := range drivers {
		driverCoord := haversine.Coord{Lat: d.Lat, Lon: d.Lon}

		_, km := haversine.Distance(userCoord, driverCoord)

		if km <= 6.0 {
			nearby = append(nearby, dto.NearbyDriverResponse{
				FirstName:  d.FirstName,
				LastName:   d.LastName,
				Plate:      d.Plate,
				DistanceKm: km,
			})
		}
	}
	sort.Slice(nearby, func(i, j int) bool {
		return nearby[i].DistanceKm < nearby[j].DistanceKm
	})
	return nearby, nil
}
