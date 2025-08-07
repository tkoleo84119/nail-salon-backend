package errors

// Error codes constants for easy reference
const (
	// AUTH - Authentication related errors
	AuthInvalidCredentials  = "AUTH_INVALID_CREDENTIALS"
	AuthTokenInvalid        = "AUTH_TOKEN_INVALID"
	AuthTokenMissing        = "AUTH_TOKEN_MISSING"
	AuthTokenFormatError    = "AUTH_TOKEN_FORMAT_ERROR"
	AuthStaffFailed         = "AUTH_STAFF_FAILED"
	AuthContextMissing      = "AUTH_CONTEXT_MISSING"
	AuthLineTokenInvalid    = "AUTH_LINE_TOKEN_INVALID"
	AuthLineTokenExpired    = "AUTH_LINE_TOKEN_EXPIRED"
	AuthRefreshTokenInvalid = "AUTH_REFRESH_TOKEN_INVALID"
	AuthPermissionDenied    = "AUTH_PERMISSION_DENIED"

	// VAL - Input validation errors
	ValJsonFormat            = "VAL_JSON_FORMAT"
	ValPathParamMissing      = "VAL_PATH_PARAM_MISSING"
	ValAllFieldsEmpty        = "VAL_ALL_FIELDS_EMPTY"
	ValTypeConversionFailed  = "VAL_TYPE_CONVERSION_FAILED"
	ValInputValidationFailed = "VAL_INPUT_VALIDATION_FAILED"

	ValFieldRequired        = "VAL_FIELD_REQUIRED"
	ValFieldStringMinLength = "VAL_FIELD_STRING_MIN_LENGTH"
	ValFieldArrayMinLength  = "VAL_FIELD_ARRAY_MIN_LENGTH"
	ValFieldMinNumber       = "VAL_FIELD_MIN_NUMBER"
	ValFieldStringMaxLength = "VAL_FIELD_STRING_MAX_LENGTH"
	ValFieldArrayMaxLength  = "VAL_FIELD_ARRAY_MAX_LENGTH"
	ValFieldMaxNumber       = "VAL_FIELD_MAX_NUMBER"
	ValFieldInvalidEmail    = "VAL_FIELD_INVALID_EMAIL"
	ValFieldNumeric         = "VAL_FIELD_NUMERIC"
	ValFieldBoolean         = "VAL_FIELD_BOOLEAN"
	ValFieldOneof           = "VAL_FIELD_ONEOF"
	ValFieldTaiwanLandline  = "VAL_FIELD_TAIWAN_LANDLINE"
	ValFieldTaiwanMobile    = "VAL_FIELD_TAIWAN_MOBILE"
	ValFieldDateFormat      = "VAL_FIELD_DATE_FORMAT"
	ValTimeConversionFailed = "VAL_TIME_CONVERSION_FAILED"
	ValFieldTimeFormat      = "VAL_FIELD_TIME_FORMAT"

	ValDuplicateWorkDate     = "VAL_DUPLICATE_WORK_DATE"
	ValTimeSlotRequired      = "VAL_TIME_SLOT_REQUIRED"
	ValEndBeforeStart        = "VAL_END_BEFORE_START"
	ValDateRangeExceed60Days = "VAL_DATE_RANGE_EXCEED_60_DAYS"

	// BOOKING - Booking operation errors
	BookingStatusNotAllowedToUpdate = "BOOKING_STATUS_NOT_ALLOWED_TO_UPDATE"
	BookingStatusNotAllowedToCancel = "BOOKING_STATUS_NOT_ALLOWED_TO_CANCEL"
	BookingNotBelongToStore         = "BOOKING_NOT_BELONG_TO_STORE"
	BookingNotFound                 = "BOOKING_NOT_FOUND"
	BookingTimeSlotNotFound         = "BOOKING_TIME_SLOT_NOT_FOUND"
	BookingTimeSlotUnavailable      = "BOOKING_TIME_SLOT_UNAVAILABLE"

	// CUSTOMER - Customer operation errors
	CustomerNotFound      = "CUSTOMER_NOT_FOUND"
	CustomerAuthNotFound  = "CUSTOMER_AUTH_NOT_FOUND"
	CustomerAlreadyExists = "CUSTOMER_ALREADY_EXISTS"

	// SCHEDULE - Schedule operation errors
	ScheduleAlreadyBookedDoNotDelete = "SCHEDULE_ALREADY_BOOKED_DO_NOT_DELETE"
	ScheduleTimeSlotInvalid          = "SCHEDULE_TIME_SLOT_INVALID"
	ScheduleNotBelongToStore         = "SCHEDULE_NOT_BELONG_TO_STORE"
	ScheduleNotBelongToStylist       = "SCHEDULE_NOT_BELONG_TO_STYLIST"
	ScheduleNotFound                 = "SCHEDULE_NOT_FOUND"
	ScheduleAlreadyExists            = "SCHEDULE_ALREADY_EXISTS"
	ScheduleEndBeforeStart           = "SCHEDULE_END_BEFORE_START"
	ScheduleDateRangeExceed31Days    = "SCHEDULE_DATE_RANGE_EXCEED_31_DAYS"
	ScheduleDuplicateWorkDateInput   = "SCHEDULE_DUPLICATE_WORK_DATE_INPUT"
	ScheduleCannotCreateBeforeToday  = "SCHEDULE_CANNOT_CREATE_BEFORE_TODAY"
	ScheduleAlreadyBookedDoNotUpdate = "SCHEDULE_ALREADY_BOOKED_DO_NOT_UPDATE"

	// SERVICE - Service operation errors
	ServiceNotActive      = "SERVICE_NOT_ACTIVE"
	ServiceNotMainService = "SERVICE_NOT_MAIN_SERVICE"
	ServiceNotAddon       = "SERVICE_NOT_ADDON"
	ServiceNotFound       = "SERVICE_NOT_FOUND"
	ServiceAlreadyExists  = "SERVICE_ALREADY_EXISTS"

	// STAFF - Staff operation errors
	StaffInvalidRole    = "STAFF_INVALID_ROLE"
	StaffStoreNotActive = "STAFF_STORE_NOT_ACTIVE"
	StaffInactive       = "STAFF_INACTIVE"
	StaffNotUpdateSelf  = "STAFF_NOT_UPDATE_SELF"
	StaffNotFound       = "STAFF_NOT_FOUND"
	StaffStoreNotFound  = "STAFF_STORE_NOT_FOUND"
	StaffAlreadyExists  = "STAFF_ALREADY_EXISTS"
	StaffEmailExists    = "STAFF_EMAIL_EXISTS"
	StaffUsernameExists = "STAFF_USERNAME_EXISTS"

	// STORE - Store operation errors
	StoreNotActive     = "STORE_NOT_ACTIVE"
	StoreNotFound      = "STORE_NOT_FOUND"
	StoreAlreadyExists = "STORE_ALREADY_EXISTS"

	// STYLIST - Stylist operation errors
	StylistNotFound = "STYLIST_NOT_FOUND"

	// TIME_SLOT - Time slot operation errors
	TimeSlotCannotUpdateSeparately          = "TIME_SLOT_CANNOT_UPDATE_SEPARATELY"
	TimeSlotNotBelongToSchedule             = "TIME_SLOT_NOT_BELONG_TO_SCHEDULE"
	TimeSlotTemplateItemNotBelongToTemplate = "TIME_SLOT_TEMPLATE_ITEM_NOT_BELONG_TO_TEMPLATE"
	TimeSlotAlreadyBookedDoNotUpdate        = "TIME_SLOT_ALREADY_BOOKED_DO_NOT_UPDATE"
	TimeSlotAlreadyBookedDoNotDelete        = "TIME_SLOT_ALREADY_BOOKED_DO_NOT_DELETE"
	TimeSlotInvalidTimeRange                = "TIME_SLOT_INVALID_TIME_RANGE"
	TimeSlotNotFound                        = "TIME_SLOT_NOT_FOUND"
	TimeSlotNotEnoughTime                   = "TIME_SLOT_NOT_ENOUGH_TIME"
	TimeSlotTemplateNotFound                = "TIME_SLOT_TEMPLATE_NOT_FOUND"
	TimeSlotTemplateItemNotFound            = "TIME_SLOT_TEMPLATE_ITEM_NOT_FOUND"
	TimeSlotConflict                        = "TIME_SLOT_CONFLICT"
	TimeSlotEndBeforeStart                  = "TIME_SLOT_END_BEFORE_START"

	// SYS - System errors
	SysInternalError      = "SYS_INTERNAL_ERROR"
	SysDatabaseError      = "SYS_DATABASE_ERROR"
	SysServiceUnavailable = "SYS_SERVICE_UNAVAILABLE"
	SysTimeout            = "SYS_TIMEOUT"
)
