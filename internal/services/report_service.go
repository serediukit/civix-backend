package services

import (
	"context"

	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/internal/repository"
	"github.com/serediukit/civix-backend/pkg/util/timeutil"
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
		return nil, err
	}

	userId := ctx.Value("user_id").(uint64)

	report := &model.Report{
		UserID:      userId,
		Location:    req.Location,
		CityID:      city.CityID,
		Description: req.Description,
		CategoryID:  req.CategoryID,
	}

	if err = s.reportRepo.CreateReport(ctx, report); err != nil {
		return nil, err
	}

	report.CreateTime = timeutil.Now()
	report.UpdateTime = timeutil.Now()
	report.CurrentStatusID = model.ReportStatusNew

	return &contracts.CreateReportResponse{
		Report: report,
	}, nil
}

func (s *reportService) GetReports(ctx context.Context, req *contracts.GetReportsRequest) (*contracts.GetReportsResponse, error) {
	location := model.Location{
		Lat: req.Lat,
		Lng: req.Lon,
	}

	city, err := s.cityRepo.GetCityByLocation(ctx, location)
	if err != nil {
		return nil, err
	}

	if req.Statuses == nil {
		req.Statuses = []model.ReportStatus{}
	}

	reports, err := s.reportRepo.GetReportsByStatuses(ctx, location, city.CityID, req.Statuses, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &contracts.GetReportsResponse{
		Reports: reports,
	}, nil
}
