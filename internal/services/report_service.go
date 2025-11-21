package services

import (
	"context"

	"github.com/pkg/errors"
	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/db"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/internal/repository"
)

type ReportService interface {
	CreateReport(ctx context.Context, req *contracts.CreateReportRequest) (*contracts.CreateReportResponse, error)
	GetReports(ctx context.Context, req *contracts.GetReportsRequest) (*contracts.GetReportsResponse, error)
}

type reportService struct {
	reportRepo repository.ReportRepository
	cityRepo   repository.CityRepository
}

func NewReportService(reportRepository repository.ReportRepository, cityRepository repository.CityRepository) ReportService {
	return &reportService{
		reportRepo: reportRepository,
		cityRepo:   cityRepository,
	}
}

func (s *reportService) CreateReport(ctx context.Context, req *contracts.CreateReportRequest) (*contracts.CreateReportResponse, error) {
	city, err := s.cityRepo.GetCityByLocation(ctx, req.Location)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			city = &model.City{CityID: repository.KyivCity}
		} else {
			return nil, err
		}
	}

	report := &model.Report{
		UserID:      req.UserID,
		Location:    req.Location,
		CityID:      city.CityID,
		Description: req.Description,
		CategoryID:  req.CategoryID,
	}

	err = s.reportRepo.CreateReport(ctx, report)

}

func (s *reportService) GetReports(ctx context.Context, req *contracts.GetReportsRequest) (*contracts.GetReportsResponse, error) {
	// TODO implement me
	panic("implement me")
}
