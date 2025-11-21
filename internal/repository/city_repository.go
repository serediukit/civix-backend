package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/serediukit/civix-backend/internal/db"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/pkg/database"
)

const (
	ErrUnknownCity = "unknown city"
	KyivCity       = "661cc9c4-9cb2-48c8-9833-2aa21fd37798"
)

type CityRepository interface {
	GetCityByLocation(ctx context.Context, location model.Location) (*model.City, error)
}

type cityRepository struct {
	store *database.Store
}

func NewCityRepository(store *database.Store) CityRepository {
	return &cityRepository{store: store}
}

func (r *cityRepository) GetCityByLocation(ctx context.Context, location model.Location) (*model.City, error) {
	sql, args, err := db.SB().
		Select(
			db.TableCitiesColumnCityID,
			db.TableCitiesColumnName,
			db.TableCitiesColumnRegion,
		).
		From(db.TableCities).
		Suffix("ORDER BY "+db.TableCitiesColumnLocation+" <-> ST_SetSRID(ST_Point(?, ?), 4326) LIMIT 1", location.Lng, location.Lat).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "Get city by location [%+v]", location)
	}

	var city model.City

	err = r.store.GetDB().
		QueryRow(ctx, sql, args...).
		Scan(
			&city.CityID,
			&city.Name,
			&city.Region,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = db.ErrNotFound
		}

		return nil, errors.Wrapf(err, "Get city by location [%+v]", location)
	}

	return &city, nil
}
