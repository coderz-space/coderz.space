package problem

import (
	"context"

	"github.com/DSAwithGautam/Coderz.space/internal/config"
	db "github.com/DSAwithGautam/Coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *db.Queries
	config  *config.Config
	pool    *pgxpool.Pool
}

func NewService(queries *db.Queries, config *config.Config, pool *pgxpool.Pool) *Service {
	return &Service{
		queries: queries,
		config:  config,
		pool:    pool,
	}
}

// Problem operations - to be implemented

// Tag operations - to be implemented

// Resource operations - to be implemented

// Helper methods

func (s *Service) GetMember(ctx context.Context, orgID pgtype.UUID, userID pgtype.UUID) (*db.OrganizationMember, error) {
	member, err := s.queries.GetOrganizationMember(ctx, db.GetOrganizationMemberParams{
		OrganizationID: orgID,
		UserID:         userID,
	})
	if err != nil {
		return nil, err
	}
	return &member, nil
}
