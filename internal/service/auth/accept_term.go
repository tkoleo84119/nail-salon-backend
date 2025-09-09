package auth

import (
	"context"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	authModel "github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type AcceptTerm struct {
	queries *dbgen.Queries
}

func NewAcceptTerm(queries *dbgen.Queries) *AcceptTerm {
	return &AcceptTerm{
		queries: queries,
	}
}

func (s *AcceptTerm) AcceptTerm(ctx context.Context, req authModel.AcceptTermRequest, customerID int64) (*authModel.AcceptTermResponse, error) {
	exists, err := s.queries.CheckCustomerTermsExistsByCustomerIDAndVersion(ctx, dbgen.CheckCustomerTermsExistsByCustomerIDAndVersionParams{
		CustomerID:   customerID,
		TermsVersion: req.TermsVersion,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check customer terms acceptance", err)
	}

	if exists {
		existingRecord, err := s.queries.GetCustomerTermsAcceptanceByCustomerIDAndVersion(ctx, dbgen.GetCustomerTermsAcceptanceByCustomerIDAndVersionParams{
			CustomerID:   customerID,
			TermsVersion: req.TermsVersion,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get existing customer terms acceptance", err)
		}

		return &authModel.AcceptTermResponse{
			ID: utils.FormatID(existingRecord.ID),
		}, nil
	}

	acceptanceID := utils.GenerateID()
	now := time.Now()
	nowPg := utils.TimePtrToPgTimestamptz(&now)

	err = s.queries.CreateCustomerTermsAcceptance(ctx, dbgen.CreateCustomerTermsAcceptanceParams{
		ID:           acceptanceID,
		CustomerID:   customerID,
		TermsVersion: req.TermsVersion,
		AcceptedAt:   nowPg,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create customer terms acceptance", err)
	}

	return &authModel.AcceptTermResponse{
		ID: utils.FormatID(acceptanceID),
	}, nil
}
