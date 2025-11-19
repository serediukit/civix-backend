package contracts

import "github.com/serediukit/civix-backend/internal/model"

type GetReportsRequest struct {
	Location model.Location       `json:"location" binding:"required"`
	Statuses []model.ReportStatus `json:"statuses"`
	PageSize uint64               `json:"pageSize" binding:"max=10"`
}

type GetReportsResponse struct {
	Reports []model.Report `json:"reports"`
}
