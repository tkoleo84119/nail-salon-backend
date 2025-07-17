package errors

// Error codes constants for easy reference
const (
	// AUTH - Authentication related errors
	AuthInvalidCredentials = "AUTH_INVALID_CREDENTIALS"
	AuthTokenExpired       = "AUTH_TOKEN_EXPIRED"
	AuthTokenInvalid       = "AUTH_TOKEN_INVALID"
	AuthTokenMissing       = "AUTH_TOKEN_MISSING"
	AuthTokenFormatError   = "AUTH_TOKEN_FORMAT_ERROR"
	AuthStaffFailed        = "AUTH_STAFF_FAILED"
	AuthContextMissing     = "AUTH_CONTEXT_MISSING"
	AuthPermissionDenied   = "AUTH_PERMISSION_DENIED"

	// USER - User operation errors
	UserInvalidRole            = "USER_INVALID_ROLE"
	UserStoreNotExist          = "USER_STORE_NOT_EXIST"
	UserStoreNotActive         = "USER_STORE_NOT_ACTIVE"
	UserInactive               = "USER_INACTIVE"
	UserNotFound               = "USER_NOT_FOUND"
	UserAlreadyExists          = "USER_ALREADY_EXISTS"
	UserEmailExists            = "USER_EMAIL_EXISTS"
	UserUsernameExists         = "USER_USERNAME_EXISTS"
	UserNotUpdateSelf          = "USER_NOT_UPDATE_SELF"
	UserStaffNotFound          = "USER_STAFF_NOT_FOUND"
	UserStoreNotFound          = "USER_STORE_NOT_FOUND"

	// STYLIST - Stylist operation errors
	StylistAlreadyExists         = "STYLIST_ALREADY_EXISTS"
	StylistNotFound              = "STYLIST_NOT_FOUND"
	StylistNotCreated            = "STYLIST_NOT_CREATED"

	// SCHEDULE - Schedule operation errors
	ScheduleAlreadyExists        = "SCHEDULE_ALREADY_EXISTS"
	ScheduleNotFound             = "SCHEDULE_NOT_FOUND"
	ScheduleTimeConflict         = "SCHEDULE_TIME_CONFLICT"

	// VAL - Input validation errors
	ValJsonFormat            = "VAL_JSON_FORMAT"
	ValInputValidationFailed = "VAL_INPUT_VALIDATION_FAILED"
	ValAllFieldsEmpty        = "VAL_ALL_FIELDS_EMPTY"

	// SYS - System errors
	SysInternalError      = "SYS_INTERNAL_ERROR"
	SysDatabaseError      = "SYS_DATABASE_ERROR"
	SysServiceUnavailable = "SYS_SERVICE_UNAVAILABLE"
	SysTimeout            = "SYS_TIMEOUT"
)
