package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/services"
	"github.com/serediukit/civix-backend/pkg/util/response"
)

type ReportController interface {
	CreateReport(ctx *gin.Context)
	GetReports(ctx *gin.Context)
}

type reportController struct {
	reportService services.ReportService
}

func NewReportController(reportService services.ReportService) ReportController {
	return &reportController{
		reportService: reportService,
	}
}

func (c *reportController) CreateReport(ctx *gin.Context) {
	var req contracts.CreateReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err)
		return
	}

	resp, err := c.reportService.CreateReport(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to create report", err)
		return
	}

	response.Created(ctx, resp)
}

func (c *reportController) GetReports(ctx *gin.Context) {
	var req contracts.GetReportsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err)
		return
	}

	resp, err := c.reportService.GetReports(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get reports", err)
		return
	}

	response.Success(ctx, resp)
}

//
// func (c *ReportController) UpdateReport(ctx *gin.Context) {
//
// }
//
// func (c *ReportController) DeleteReport(ctx *gin.Context) {
//
// }
