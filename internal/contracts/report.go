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
	Lat      float64              `form:"lat" binding:"required"`
	Lon      float64              `form:"lon" binding:"required"`
	Statuses []model.ReportStatus `form:"statuses"`
	PageSize uint64               `form:"page_size" binding:"max=100"`
}

type GetReportsResponse struct {
	Reports []*model.Report `json:"reports"`
}
