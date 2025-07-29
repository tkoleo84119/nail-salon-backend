package errors

// Error codes constants for easy reference
const (
	// AUTH - Authentication related errors
	AuthInvalidCredentials  = "AUTH_INVALID_CREDENTIALS"
	AuthTokenExpired        = "AUTH_TOKEN_EXPIRED"
	AuthTokenInvalid        = "AUTH_TOKEN_INVALID"
	AuthTokenMissing        = "AUTH_TOKEN_MISSING"
	AuthTokenFormatError    = "AUTH_TOKEN_FORMAT_ERROR"
	AuthStaffFailed         = "AUTH_STAFF_FAILED"
	AuthContextMissing      = "AUTH_CONTEXT_MISSING"
	AuthPermissionDenied    = "AUTH_PERMISSION_DENIED"
	AuthLineTokenInvalid    = "AUTH_LINE_TOKEN_INVALID"
	AuthLineTokenExpired    = "AUTH_LINE_TOKEN_EXPIRED"
	AuthRefreshTokenInvalid = "AUTH_REFRESH_TOKEN_INVALID"

	// BOOKING - Booking operation errors
	BookingNotFound                 = "BOOKING_NOT_FOUND"
	BookingStatusNotAllowedToUpdate = "BOOKING_STATUS_NOT_ALLOWED_TO_UPDATE"
	BookingStatusNotAllowedToCancel = "BOOKING_STATUS_NOT_ALLOWED_TO_CANCEL"
	BookingTimeSlotNotFound         = "TIME_SLOT_NOT_FOUND"
	BookingTimeSlotUnavailable      = "TIME_SLOT_UNAVAILABLE"

	// CUSTOMER - Customer operation errors
	CustomerNotFound      = "CUSTOMER_NOT_FOUND"
	CustomerAuthNotFound  = "CUSTOMER_AUTH_NOT_FOUND"
	CustomerAlreadyExists = "CUSTOMER_ALREADY_EXISTS"

	// SCHEDULE - Schedule operation errors
	ScheduleAlreadyExists            = "SCHEDULE_ALREADY_EXISTS"
	ScheduleNotFound                 = "SCHEDULE_NOT_FOUND"
	ScheduleAlreadyBookedDoNotDelete = "SCHEDULE_ALREADY_BOOKED_DO_NOT_DELETE"
	ScheduleNotBelongToStore         = "SCHEDULE_NOT_BELONG_TO_STORE"
	ScheduleNotBelongToStylist       = "SCHEDULE_NOT_BELONG_TO_STYLIST"

	// SERVICE - Service operation errors
	ServiceNotActive      = "SERVICE_NOT_ACTIVE"
	ServiceNotMainService = "SERVICE_NOT_MAIN_SERVICE"
	ServiceNotAddon       = "SERVICE_NOT_ADDON"
	ServiceNotFound       = "SERVICE_NOT_FOUND"
	ServiceAlreadyExists  = "SERVICE_ALREADY_EXISTS"

	// STORE - Store operation errors
	StoreNotFound      = "STORE_NOT_FOUND"
	StoreNotActive     = "STORE_NOT_ACTIVE"
	StoreAlreadyExists = "STORE_ALREADY_EXISTS"

	// STYLIST - Stylist operation errors
	StylistAlreadyExists = "STYLIST_ALREADY_EXISTS"
	StylistNotFound      = "STYLIST_NOT_FOUND"
	StylistNotCreated    = "STYLIST_NOT_CREATED"

	// TIME_SLOT - Time slot operation errors
	TimeSlotCannotUpdateSeparately          = "TIME_SLOT_CANNOT_UPDATE_SEPARATELY"
	TimeSlotNotBelongToSchedule             = "TIME_SLOT_NOT_BELONG_TO_SCHEDULE"
	TimeSlotTemplateItemNotBelongToTemplate = "TIME_SLOT_TEMPLATE_ITEM_NOT_BELONG_TO_TEMPLATE"
	TimeSlotAlreadyBookedDoNotUpdate        = "TIME_SLOT_ALREADY_BOOKED_DO_NOT_UPDATE"
	TimeSlotAlreadyBookedDoNotDelete        = "TIME_SLOT_ALREADY_BOOKED_DO_NOT_DELETE"
	TimeSlotInvalidTimeRange                = "TIME_SLOT_INVALID_TIME_RANGE"
	TimeSlotConflict                        = "TIME_SLOT_CONFLICT"
	TimeSlotNotFound                        = "TIME_SLOT_NOT_FOUND"
	TimeSlotNotEnoughTime                   = "TIME_SLOT_NOT_ENOUGH_TIME"
	TimeSlotTemplateNotFound                = "TIME_SLOT_TEMPLATE_NOT_FOUND"
	TimeSlotTemplateItemNotFound            = "TIME_SLOT_TEMPLATE_ITEM_NOT_FOUND"

	// USER - User operation errors
	UserInvalidRole    = "USER_INVALID_ROLE"
	UserStoreNotActive = "USER_STORE_NOT_ACTIVE"
	UserInactive       = "USER_INACTIVE"
	UserNotUpdateSelf  = "USER_NOT_UPDATE_SELF"
	UserNotFound       = "USER_NOT_FOUND"
	UserStoreNotFound  = "USER_STORE_NOT_FOUND"
	UserAlreadyExists  = "USER_ALREADY_EXISTS"
	UserEmailExists    = "USER_EMAIL_EXISTS"
	UserUsernameExists = "USER_USERNAME_EXISTS"

	// VAL - Input validation errors
	ValJsonFormat            = "VAL_JSON_FORMAT"
	ValInputValidationFailed = "VAL_INPUT_VALIDATION_FAILED"
	ValAllFieldsEmpty        = "VAL_ALL_FIELDS_EMPTY"
	ValDateFormatInvalid     = "VAL_DATE_FORMAT_INVALID"
	ValDuplicateWorkDate     = "VAL_DUPLICATE_WORK_DATE"
	ValTimeSlotRequired      = "VAL_TIME_SLOT_REQUIRED"
	ValEndBeforeStart        = "VAL_END_BEFORE_START"
	ValDateRangeExceed60Days = "VAL_DATE_RANGE_EXCEED_60_DAYS"

	// SYS - System errors
	SysInternalError      = "SYS_INTERNAL_ERROR"
	SysDatabaseError      = "SYS_DATABASE_ERROR"
	SysServiceUnavailable = "SYS_SERVICE_UNAVAILABLE"
	SysTimeout            = "SYS_TIMEOUT"
)
