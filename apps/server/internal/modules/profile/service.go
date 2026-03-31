package profile

import (
	"context"

	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
	db "github.com/DSAwithGautam/Coderz.space/internal/db/sqlc"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

type ProfileDTO struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Solved    int32  `json:"solved"`
	JoinedAt  string `json:"joinedAt"`
}

func (s *Service) GetMenteeProfile(ctx context.Context, userIDStr string) (*ProfileDTO, error) {
	userID, err := utils.StringToUUID(userIDStr)
	if err != nil {
		return nil, err
	}

	r, err := s.queries.WebGetProfileStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &ProfileDTO{
		FirstName: r.FirstName,
		LastName:  "", // name wasn't splitted in DB
		Solved:    int32(r.Solved),
		JoinedAt:  r.JoinedAt.Time.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *Service) GetLeaderboard(ctx context.Context) ([]ProfileDTO, error) {
	rows, err := s.queries.WebGetLeaderboard(ctx)
	if err != nil {
		return nil, err
	}

	var res []ProfileDTO
	for _, r := range rows {
		res = append(res, ProfileDTO{
			ID:        utils.UUIDToString(r.ID),
			FirstName: r.FirstName,
			LastName:  r.LastName,
			Solved:    r.Solved,
		})
	}
	// never return nil array
	if res == nil {
		res = []ProfileDTO{}
	}
	return res, nil
}
