package db

const (
	TableUsers   = "users"
	TableReports = "reports"
	TableCities  = "cities"
)

const (
	TableUsersColumnUserID       = "user_id"
	TableUsersColumnEmail        = "email"
	TableUsersColumnPasswordHash = "password_hash"
	TableUsersColumnName         = "name"
	TableUsersColumnSurname      = "surname"
	TableUsersColumnPhoneNumber  = "phone_number"
	TableUsersColumnAvatarUrl    = "avatar_url"
	TableUsersColumnRegCityID    = "reg_city_id"
	TableUsersColumnRegTime      = "reg_time"
	TableUsersColumnUpdTime      = "upd_time"
	TableUsersColumnDelTime      = "del_time"
)

const (
	TableReportsColumnReportID  = "report_id"
	TableReportsColumnUserID    = "user_id"
	TableReportsCreateTime      = "create_time"
	TableReportsUpdateTime      = "update_time"
	TableReportsLocation        = "location"
	TableReportsCityID          = "city_id"
	TableReportsDescription     = "description"
	TableReportsCategoryID      = "category_id"
	TableReportsCurrentStatusID = "current_status_id"
)

const (
	TableCitiesColumnCityID   = "city_id"
	TableCitiesColumnName     = "name"
	TableCitiesColumnRegion   = "region"
	TableCitiesColumnLocation = "location"
)
