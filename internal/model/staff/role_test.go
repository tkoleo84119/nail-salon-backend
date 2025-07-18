package staff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{
			name:     "valid SUPER_ADMIN role",
			role:     RoleSuperAdmin,
			expected: true,
		},
		{
			name:     "valid ADMIN role",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "valid MANAGER role",
			role:     RoleManager,
			expected: true,
		},
		{
			name:     "valid STYLIST role",
			role:     RoleStylist,
			expected: true,
		},
		{
			name:     "invalid role",
			role:     "INVALID_ROLE",
			expected: false,
		},
		{
			name:     "empty role",
			role:     "",
			expected: false,
		},
		{
			name:     "lowercase role",
			role:     "admin",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidRole(tt.role)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAllRoles(t *testing.T) {
	roles := GetAllRoles()

	expectedRoles := []string{
		RoleSuperAdmin,
		RoleAdmin,
		RoleManager,
		RoleStylist,
	}

	assert.Equal(t, expectedRoles, roles)
	assert.Len(t, roles, 4)
}

func TestRoleConstants(t *testing.T) {
	assert.Equal(t, "SUPER_ADMIN", RoleSuperAdmin)
	assert.Equal(t, "ADMIN", RoleAdmin)
	assert.Equal(t, "MANAGER", RoleManager)
	assert.Equal(t, "STYLIST", RoleStylist)
}
