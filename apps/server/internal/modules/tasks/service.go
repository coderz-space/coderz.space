package tasks

import (
	"context"
	"errors"

	"github.com/DSAwithGautam/Coderz.space/internal/common/utils"
	db "github.com/DSAwithGautam/Coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) ListMenteeQuestions(ctx context.Context, userIDStr string) ([]QuestionDTO, error) {
	userID, err := utils.StringToUUID(userIDStr)
	if err != nil {
		return nil, err
	}

	rows, err := s.queries.WebListMenteeQuestions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var res []QuestionDTO
	for _, r := range rows {
		assignedAtStr := ""
		if r.AssignedAt.Valid {
			assignedAtStr = r.AssignedAt.Time.Format("2006-01-02T15:04:05Z")
		}
		completedAtStr := ""
		if r.CompletedAt.Valid {
			completedAtStr = r.CompletedAt.Time.Format("2006-01-02T15:04:05Z")
		}

		// map DB enums correctly
		res = append(res, QuestionDTO{
			ID:             utils.UUIDToString(r.ID),
			Title:          r.Title,
			Description:    r.Description.String,
			Difficulty:     string(r.Difficulty),
			Topic:          "", // topic not explicitly saved in problem table directly yet
			Status:         string(r.AssignmentStatus),
			ProgressStatus: string(r.ProgressStatus),
			AssignedAt:     assignedAtStr,
			CompletedAt:    completedAtStr,
			SolutionUrl:    r.SolutionUrl.String,
			Solution:       r.Solution.String,
			Resources:      r.Resources,
		})
	}
	// return empty array instead of null
	if res == nil {
		res = []QuestionDTO{}
	}
	return res, nil
}

func (s *Service) AssignQuestionToMentee(ctx context.Context, userIDStr string, req AssignQuestionRequest) (*QuestionDTO, error) {
	// Not fully implementing dynamic Problem creation logic since MVP assumes 
	// questions are already populated or mentor will use existing problems.
	// For fallback test purpose we fail over:
	return nil, errors.New("Assigning tasks requires full bootcamp assignment mapping flow - stubbed for MVP")
}

func (s *Service) GetMenteeQuestion(ctx context.Context, userIDStr, questionIDStr string) (*QuestionDTO, error) {
	userID, err := utils.StringToUUID(userIDStr)
	if err != nil {
		return nil, err
	}
	questionID, err := utils.StringToUUID(questionIDStr)
	if err != nil {
		return nil, err
	}

	r, err := s.queries.WebGetMenteeQuestion(ctx, db.WebGetMenteeQuestionParams{
		UserID: userID,
		ID:     questionID,
	})
	if err != nil {
		return nil, err
	}

	assignedAtStr := ""
	if r.AssignedAt.Valid {
		assignedAtStr = r.AssignedAt.Time.Format("2006-01-02T15:04:05Z")
	}
	completedAtStr := ""
	if r.CompletedAt.Valid {
		completedAtStr = r.CompletedAt.Time.Format("2006-01-02T15:04:05Z")
	}

	return &QuestionDTO{
		ID:             utils.UUIDToString(r.ID),
		Title:          r.Title,
		Description:    r.Description.String,
		Difficulty:     string(r.Difficulty),
		Status:         string(r.AssignmentStatus),
		ProgressStatus: string(r.ProgressStatus),
		AssignedAt:     assignedAtStr,
		CompletedAt:    completedAtStr,
		SolutionUrl:    r.SolutionUrl.String,
		Solution:       r.Solution.String,
	}, nil
}

func (s *Service) UpdateQuestionProgress(ctx context.Context, userIDStr, questionIDStr, progress string) error {
	userID, err := utils.StringToUUID(userIDStr)
	if err != nil {
		return err
	}
	questionID, err := utils.StringToUUID(questionIDStr)
	if err != nil {
		return err
	}

	return s.queries.WebUpdateQuestionProgress(ctx, db.WebUpdateQuestionProgressParams{
		UserID: userID,
		ProblemID: questionID,
		Status: db.AssignmentProblemStatus(progress),
	})
}

func (s *Service) UpdateQuestionDetails(ctx context.Context, userIDStr, questionIDStr string, req UpdateDetailsRequest) error {
	userID, err := utils.StringToUUID(userIDStr)
	if err != nil {
		return err
	}
	questionID, err := utils.StringToUUID(questionIDStr)
	if err != nil {
		return err
	}

	return s.queries.WebUpdateQuestionDetails(ctx, db.WebUpdateQuestionDetailsParams{
		UserID: userID,
		ProblemID: questionID,
		Notes: pgtype.Text{String: req.Solution, Valid: req.Solution != ""},
		SolutionLink: pgtype.Text{String: req.Resources, Valid: req.Resources != ""}, // Mapping resources string to solution_link as arbitrary text
	})
}
