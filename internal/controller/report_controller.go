package controller

import "github.com/gin-gonic/gin"

type ReportController struct {
	reportService report.ReportService
}

func NewReportController(reportService *report.ReportService) *ReportController {
	return &ReportController{
		reportService: reportService,
	}
}

func (c *ReportController) GetReports(ctx *gin.Context) {

}

func (c *ReportController) CreateReport(ctx *gin.Context) {

}

func (c *ReportController) UpdateReport(ctx *gin.Context) {

}

func (c *ReportController) DeleteReport(ctx *gin.Context) {

}
