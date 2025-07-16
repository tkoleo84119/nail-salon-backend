package stylist

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)

func TestCreateStylistService_CreateStylist_Success(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStylistService(mockQuerier)

	ctx := context.Background()
	staffUserID := int64(12345)
	req := stylist.CreateStylistRequest{
		StylistName:  "Jane 美甲師",
		GoodAtShapes: []string{"方形", "圓形"},
		GoodAtColors: []string{"裸色系", "粉嫩系"},
		GoodAtStyles: []string{"手繪", "簡約"},
		IsIntrovert:  boolPtr(false),
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleAdmin,
	}, nil)

	// Mock stylist existence check
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(false, nil)

	// Mock stylist creation
	mockQuerier.On("CreateStylist", ctx, mock.AnythingOfType("dbgen.CreateStylistParams")).Return(dbgen.Stylist{
		ID:           18000000001,
		StaffUserID:  pgtype.Int8{Int64: staffUserID, Valid: true},
		Name:         pgtype.Text{String: req.StylistName, Valid: true},
		GoodAtShapes: req.GoodAtShapes,
		GoodAtColors: req.GoodAtColors,
		GoodAtStyles: req.GoodAtStyles,
		IsIntrovert:  pgtype.Bool{Bool: false, Valid: true},
	}, nil)

	response, err := service.CreateStylist(ctx, req, staffUserID)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "18000000001", response.ID)
	assert.Equal(t, "12345", response.StaffUserID)
	assert.Equal(t, "Jane 美甲師", response.StylistName)
	assert.Equal(t, []string{"方形", "圓形"}, response.GoodAtShapes)
	assert.Equal(t, []string{"裸色系", "粉嫩系"}, response.GoodAtColors)
	assert.Equal(t, []string{"手繪", "簡約"}, response.GoodAtStyles)
	assert.Equal(t, false, response.IsIntrovert)

	mockQuerier.AssertExpectations(t)
}

func TestCreateStylistService_CreateStylist_SuperAdminForbidden(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStylistService(mockQuerier)

	ctx := context.Background()
	staffUserID := int64(12345)
	req := stylist.CreateStylistRequest{
		StylistName: "Jane 美甲師",
	}

	// Mock staff user lookup - SUPER_ADMIN role
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleSuperAdmin,
	}, nil)

	response, err := service.CreateStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateStylistService_CreateStylist_AlreadyExists(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStylistService(mockQuerier)

	ctx := context.Background()
	staffUserID := int64(12345)
	req := stylist.CreateStylistRequest{
		StylistName: "Jane 美甲師",
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleAdmin,
	}, nil)

	// Mock stylist existence check - already exists
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(true, nil)

	response, err := service.CreateStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.StylistAlreadyExists, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateStylistService_CreateStylist_StaffUserNotFound(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStylistService(mockQuerier)

	ctx := context.Background()
	staffUserID := int64(12345)
	req := stylist.CreateStylistRequest{
		StylistName: "Jane 美甲師",
	}

	// Mock staff user lookup - not found
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{}, errors.New("staff user not found"))

	response, err := service.CreateStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthStaffFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateStylistService_CreateStylist_DatabaseError(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStylistService(mockQuerier)

	ctx := context.Background()
	staffUserID := int64(12345)
	req := stylist.CreateStylistRequest{
		StylistName: "Jane 美甲師",
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleAdmin,
	}, nil)

	// Mock stylist existence check
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(false, nil)

	// Mock stylist creation - database error
	mockQuerier.On("CreateStylist", ctx, mock.AnythingOfType("dbgen.CreateStylistParams")).Return(dbgen.Stylist{}, errors.New("database error"))

	response, err := service.CreateStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateStylistService_CreateStylist_WithDefaults(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStylistService(mockQuerier)

	ctx := context.Background()
	staffUserID := int64(12345)
	req := stylist.CreateStylistRequest{
		StylistName: "Jane 美甲師",
		// No optional fields provided
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleAdmin,
	}, nil)

	// Mock stylist existence check
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(false, nil)

	// Mock stylist creation
	mockQuerier.On("CreateStylist", ctx, mock.AnythingOfType("dbgen.CreateStylistParams")).Return(dbgen.Stylist{
		ID:           18000000001,
		StaffUserID:  pgtype.Int8{Int64: staffUserID, Valid: true},
		Name:         pgtype.Text{String: req.StylistName, Valid: true},
		GoodAtShapes: []string{},
		GoodAtColors: []string{},
		GoodAtStyles: []string{},
		IsIntrovert:  pgtype.Bool{Bool: false, Valid: true},
	}, nil)

	response, err := service.CreateStylist(ctx, req, staffUserID)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "18000000001", response.ID)
	assert.Equal(t, "12345", response.StaffUserID)
	assert.Equal(t, "Jane 美甲師", response.StylistName)
	assert.Equal(t, []string{}, response.GoodAtShapes)
	assert.Equal(t, []string{}, response.GoodAtColors)
	assert.Equal(t, []string{}, response.GoodAtStyles)
	assert.Equal(t, false, response.IsIntrovert)

	mockQuerier.AssertExpectations(t)
}

// Helper function to create boolean pointer
func boolPtr(b bool) *bool {
	return &b
}
