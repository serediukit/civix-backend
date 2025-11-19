package services

import (
	"context"

	"github.com/serediukit/civix-backend/internal/contracts"
)

type ReportService interface {
	CreateReport(ctx context.Context, req *contracts.CreateReportRequest) (*contracts.CreateReportResponse, error)
	GetReports(ctx context.Context, req *contracts.GetReportsRequest) (*contracts.GetReportsResponse, error)
}
