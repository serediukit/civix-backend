package repository

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/serediukit/civix-backend/internal/db"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/pkg/database"
)

type ReportRepository interface {
	CreateReport(ctx context.Context, req *model.Report) error
	GetReportsByStatuses(ctx context.Context, location model.Location, cityID string, statuses []model.ReportStatus, pageSize uint64) ([]*model.Report, error)
}

type reportRepository struct {
	store *database.Store
}

func NewReportRepository(store *database.Store) ReportRepository {
	return &reportRepository{store: store}
}

func (r *reportRepository) CreateReport(ctx context.Context, report *model.Report) error {
	sql, args, err := db.SB().
		Insert(db.TableReports).
		Columns(
			db.TableReportsColumnUserID,
			db.TableReportsLocation,
			db.TableReportsCityID,
			db.TableReportsDescription,
			db.TableReportsCategoryID,
		).
		Values(
			report.UserID,
			squirrel.Expr("ST_SetSRID(ST_Point(?, ?), 4326)", report.Location.Lng, report.Location.Lat),
			report.CityID,
			report.Description,
			report.CategoryID,
		).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, "Create report [%+v] ToSQL: %s, %+v", report, sql, args)
	}

	_, err = r.store.GetDB().Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrapf(err, "Create report [%+v] Exec: %s, %+v", report, sql, args)
	}

	return nil
}

func (r *reportRepository) GetReportsByStatuses(ctx context.Context, location model.Location, cityID string, statuses []model.ReportStatus, pageSize uint64) ([]*model.Report, error) {
	sb := db.SB().
		Select(
			db.TableReportsColumnReportID,
			db.TableReportsColumnUserID,
			db.TableReportsCreateTime,
			db.TableReportsUpdateTime,
			"ST_X("+db.TableReportsLocation+") as lon",
			"ST_Y("+db.TableReportsLocation+") as lat",
			db.TableReportsCityID,
			db.TableReportsDescription,
			db.TableReportsCategoryID,
			db.TableReportsCurrentStatusID,
		).
		From(db.TableReports).
		Where(squirrel.Eq{db.TableReportsCityID: cityID})

	if len(statuses) > 0 {
		sb = sb.Where(squirrel.Eq{db.TableReportsCurrentStatusID: statuses})
	}

	sql, args, err := sb.
		Suffix("ORDER BY "+db.TableReportsLocation+" <-> ST_SetSRID(ST_Point(?, ?), 4326) LIMIT ?", location.Lng, location.Lat, pageSize).
		ToSql()

	log.Printf("| GetReportsByStatuses |\nSQL: %s\nARGS: %+v\n", sql, args)

	rows, err := r.store.GetDB().Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = db.ErrNotFound
		}

		return nil, errors.Wrapf(err, "Get reports by statuses %+v] for city_id [%s] and location [%+v]", statuses, cityID, location)
	}

	reports := make([]*model.Report, 0)

	for rows.Next() {
		report := &model.Report{}

		if err = rows.Scan(
			&report.ReportID,
			&report.UserID,
			&report.CreateTime,
			&report.UpdateTime,
			&report.Location.Lng,
			&report.Location.Lat,
			&report.CityID,
			&report.Description,
			&report.CategoryID,
			&report.CurrentStatusID,
		); err != nil {
			return nil, errors.Wrapf(err, "Get reports by statuses [%+v] for city_id [%s] and location [%+v]", statuses, cityID, location)
		}

		reports = append(reports, report)
	}

	return reports, nil
}
