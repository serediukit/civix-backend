package model

type ReportStatus int8

const (
	ReportStatusNew ReportStatus = iota
	ReportStatusInProgress
	ReportStatusCompleted
)

type ReportCategory int8

const (
	ReportCategoryUnknown ReportCategory = iota
	ReportCategoryRoad
	ReportCategorySideway
	ReportCategoryElectric
	ReportCategoryWater
	ReportCategoryGas
	ReportCategoryHeat
)

type Report struct {
	ReportID        string         `json:"report_id"`
	UserID          uint64         `json:"user_id"`
	CreateTime      string         `json:"create_time"`
	UpdateTime      string         `json:"update_time"`
	Location        Location       `json:"location"`
	CityID          string         `json:"city_id"`
	Description     string         `json:"description"`
	CategoryID      ReportCategory `json:"category_id"`
	CurrentStatusID ReportStatus   `json:"current_status_id"`
}
