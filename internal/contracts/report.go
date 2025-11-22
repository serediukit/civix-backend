package contracts

import "github.com/serediukit/civix-backend/internal/model"

type CreateReportRequest struct {
	Location    model.Location       `json:"location" binding:"required"`
	Description string               `json:"description"`
	CategoryID  model.ReportCategory `json:"category_id"`
}

type CreateReportResponse struct {
	Report *model.Report `json:"report"`
}

type GetReportsRequest struct {
	Location model.Location       `json:"location" binding:"required"`
	Statuses []model.ReportStatus `json:"statuses"`
	PageSize uint64               `json:"page_size" binding:"max=10"`
}

type GetReportsResponse struct {
	Reports []*model.Report `json:"reports"`
}
