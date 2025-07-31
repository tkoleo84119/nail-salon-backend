package common

// Role constants for staff users
const (
	RoleSuperAdmin = "SUPER_ADMIN"
	RoleAdmin      = "ADMIN"
	RoleManager    = "MANAGER"
	RoleStylist    = "STYLIST"
)

// IsValidRole checks if the given role is valid
func IsValidRole(role string) bool {
	switch role {
	case RoleSuperAdmin, RoleAdmin, RoleManager, RoleStylist:
		return true
	default:
		return false
	}
}

// GetAllRoles returns all valid roles
func GetAllRoles() []string {
	return []string{
		RoleSuperAdmin,
		RoleAdmin,
		RoleManager,
		RoleStylist,
	}
}