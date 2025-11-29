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
	TableReportsColumnReportID        = "report_id"
	TableReportsColumnUserID          = "user_id"
	TableReportsColumnCreateTime      = "create_time"
	TableReportsColumnUpdateTime      = "update_time"
	TableReportsColumnLocation        = "location"
	TableReportsColumnCityID          = "city_id"
	TableReportsColumnDescription     = "description"
	TableReportsColumnCategoryID      = "category_id"
	TableReportsColumnCurrentStatusID = "current_status_id"
	TableReportsColumnPhotoURL        = "photo_url"
)

const (
	TableCitiesColumnCityID   = "city_id"
	TableCitiesColumnName     = "name"
	TableCitiesColumnRegion   = "region"
	TableCitiesColumnLocation = "location"
)
