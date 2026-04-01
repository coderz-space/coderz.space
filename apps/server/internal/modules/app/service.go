package app

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	db "github.com/coderz-space/coderz.space/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var usernamePattern = regexp.MustCompile(`^[a-z0-9_]+$`)

type Service struct {
	pool *pgxpool.Pool
}

type resolvedContext struct {
	User           UserData
	UserID         string
	Organization   *OrganizationData
	Bootcamp       *BootcampData
	MemberID       string
	OrgRole        string
	EnrollmentID   string
	EnrollmentRole string
	AssignedSheet  string
	Role           string
	AccountStatus  string
}

type menteeRecord struct {
	EnrollmentID  string
	MemberID      string
	UserID        string
	FirstName     string
	LastName      string
	Username      string
	Email         string
	AssignedSheet string
	EnrolledAt    time.Time
}

type questionRow struct {
	ID             string
	AssignmentID   string
	TargetUsername string
	Title          string
	Description    string
	Difficulty     string
	ExternalLink   string
	AppProgress    string
	LegacyStatus   string
	Notes          string
	Resources      string
	AssignedAt     time.Time
	CompletedAt    pgtype.Timestamptz
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) GetContext(ctx context.Context, userID string) (*ContextData, error) {
	resolved, err := s.resolveContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	data := &ContextData{
		Role:          resolved.Role,
		AccountStatus: resolved.AccountStatus,
		User:          resolved.User,
		Organization:  resolved.Organization,
		Bootcamp:      resolved.Bootcamp,
	}
	if resolved.EnrollmentID != "" || resolved.AssignedSheet != "" {
		data.Enrollment = &EnrollmentData{
			ID:            resolved.EnrollmentID,
			AssignedSheet: resolved.AssignedSheet,
		}
	}

	return data, nil
}

func (s *Service) MenteeSignup(ctx context.Context, req MenteeSignupRequest) (*MenteeSignupData, error) {
	username := normalizeUsername(req.Username)
	if !usernamePattern.MatchString(username) {
		return nil, errors.New("INVALID_USERNAME")
	}
	if !validatePasswordComplexity(req.Password) {
		return nil, errors.New("PASSWORD_MUST_CONTAIN_LETTER_AND_NUMBER")
	}

	defaultOrg, defaultBootcamp, err := s.getDefaultSignupContext(ctx)
	if err != nil {
		return nil, err
	}

	fullName := strings.TrimSpace(strings.TrimSpace(req.FirstName) + " " + strings.TrimSpace(req.LastName))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var emailExists bool
	if err := tx.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM coderz.users
			WHERE LOWER(COALESCE(email, '')) = LOWER($1)
		)
	`, req.Email).Scan(&emailExists); err != nil {
		return nil, err
	}
	if emailExists {
		return nil, errors.New("EMAIL_ALREADY_EXISTS")
	}

	var usernameExists bool
	if err := tx.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM coderz.users
			WHERE LOWER(username) = LOWER($1)
		)
	`, username).Scan(&usernameExists); err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, errors.New("USERNAME_ALREADY_EXISTS")
	}

	var userIDValue string
	if err := tx.QueryRow(ctx, `
		INSERT INTO coderz.users (
			name,
			email,
			password_hash,
			role,
			username
		) VALUES (
			$1,
			$2,
			$3,
			'user',
			$4
		)
		RETURNING id::text
	`, fullName, req.Email, string(hashedPassword), username).Scan(&userIDValue); err != nil {
		return nil, err
	}

	var requestID string
	if err := tx.QueryRow(ctx, `
		INSERT INTO coderz.mentee_requests (
			user_id,
			organization_id,
			bootcamp_id,
			status
		) VALUES (
			$1,
			$2,
			$3,
			'pending'
		)
		RETURNING id::text
	`, userIDValue, defaultOrg.ID, defaultBootcamp.ID).Scan(&requestID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &MenteeSignupData{
		RequestID: requestID,
		Status:    "pending",
		Username:  username,
		Email:     req.Email,
	}, nil
}

func (s *Service) ListMenteeRequests(ctx context.Context, userID string) ([]MenteeRequestData, error) {
	resolved, err := s.resolveMentorContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	rows, err := s.pool.Query(ctx, `
		SELECT
			mr.id::text,
			u.name,
			COALESCE(u.username, ''),
			COALESCE(u.email, ''),
			mr.created_at,
			mr.status,
			COALESCE(mr.sheet_key, '')
		FROM coderz.mentee_requests mr
		JOIN coderz.users u ON u.id = mr.user_id
		WHERE mr.bootcamp_id = $1
		ORDER BY
			CASE WHEN mr.status = 'pending' THEN 0 ELSE 1 END,
			mr.created_at DESC
	`, resolved.Bootcamp.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]MenteeRequestData, 0)
	for rows.Next() {
		var (
			requestID     string
			fullName      string
			username      string
			email         string
			signedUpAt    time.Time
			status        string
			assignedSheet string
		)
		if err := rows.Scan(&requestID, &fullName, &username, &email, &signedUpAt, &status, &assignedSheet); err != nil {
			return nil, err
		}

		firstName, lastName := splitName(fullName)
		requests = append(requests, MenteeRequestData{
			ID:            requestID,
			FirstName:     firstName,
			LastName:      lastName,
			Username:      username,
			Email:         email,
			SignedUpAt:    signedUpAt.Format(time.RFC3339),
			Status:        status,
			AssignedSheet: assignedSheet,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (s *Service) ReviewMenteeRequest(ctx context.Context, userID, requestID string, req ReviewMenteeRequest) (*MenteeRequestData, error) {
	resolved, err := s.resolveMentorContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	if req.Status == "approved" && req.SheetKey == "" {
		return nil, errors.New("SHEET_REQUIRED")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var (
		targetUserID string
		orgID        string
		bootcampID   string
		fullName     string
		username     string
		email        string
		createdAt    time.Time
	)
	if err := tx.QueryRow(ctx, `
		SELECT
			mr.user_id::text,
			mr.organization_id::text,
			mr.bootcamp_id::text,
			u.name,
			COALESCE(u.username, ''),
			COALESCE(u.email, ''),
			mr.created_at
		FROM coderz.mentee_requests mr
		JOIN coderz.users u ON u.id = mr.user_id
		WHERE mr.id = $1
		  AND mr.bootcamp_id = $2
	`, requestID, resolved.Bootcamp.ID).Scan(
		&targetUserID,
		&orgID,
		&bootcampID,
		&fullName,
		&username,
		&email,
		&createdAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("REQUEST_NOT_FOUND")
		}
		return nil, err
	}

	if req.Status == "approved" {
		memberID, err := s.ensureOrganizationMember(ctx, tx, orgID, targetUserID)
		if err != nil {
			return nil, err
		}
		if err := s.ensureBootcampEnrollment(ctx, tx, bootcampID, memberID, req.SheetKey); err != nil {
			return nil, err
		}
	}

	if _, err := tx.Exec(ctx, `
		UPDATE coderz.mentee_requests
		SET
			status = $2,
			sheet_key = NULLIF($3, ''),
			reviewed_by = $4,
			reviewed_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, requestID, req.Status, req.SheetKey, resolved.MemberID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	firstName, lastName := splitName(fullName)
	return &MenteeRequestData{
		ID:            requestID,
		FirstName:     firstName,
		LastName:      lastName,
		Username:      username,
		Email:         email,
		SignedUpAt:    createdAt.Format(time.RFC3339),
		Status:        req.Status,
		AssignedSheet: req.SheetKey,
	}, nil
}

func (s *Service) ListSheets() []SheetData {
	return listSheets()
}

func (s *Service) GetDayAssignments(ctx context.Context, userID, day string) (*DayAssignmentsData, error) {
	resolved, err := s.resolveMentorContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	normalizedDay := normalizeDay(day)
	rows, err := s.pool.Query(ctx, `
		SELECT
			u.name,
			COALESCE(u.username, ''),
			COALESCE(u.email, ''),
			COALESCE(be.assigned_sheet_key, ''),
			(mda.id IS NOT NULL) AS assigned
		FROM coderz.bootcamp_enrollments be
		JOIN coderz.organization_members om ON om.id = be.organization_member_id
		JOIN coderz.users u ON u.id = om.user_id
		LEFT JOIN coderz.mentee_day_assignments mda
			ON mda.bootcamp_enrollment_id = be.id
			AND mda.weekday = $2
		WHERE be.bootcamp_id = $1
		  AND be.role = 'mentee'
		  AND be.status = 'active'
		ORDER BY u.name ASC
	`, resolved.Bootcamp.ID, normalizedDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mentees := make([]DayAssignmentMenteeData, 0)
	for rows.Next() {
		var (
			fullName      string
			username      string
			email         string
			assignedSheet string
			assigned      bool
		)
		if err := rows.Scan(&fullName, &username, &email, &assignedSheet, &assigned); err != nil {
			return nil, err
		}

		firstName, lastName := splitName(fullName)
		mentees = append(mentees, DayAssignmentMenteeData{
			FirstName:     firstName,
			LastName:      lastName,
			Username:      username,
			Email:         email,
			Assigned:      assigned,
			AssignedSheet: assignedSheet,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &DayAssignmentsData{
		Day:     normalizedDay,
		Mentees: mentees,
	}, nil
}

func (s *Service) UpdateDayAssignments(ctx context.Context, userID, day string, req UpdateDayAssignmentsRequest) (*DayAssignmentsData, error) {
	resolved, err := s.resolveMentorContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	normalizedDay := normalizeDay(day)
	targets := dedupeLower(req.Usernames)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	enrollmentMap, err := s.listMenteeEnrollmentMap(ctx, tx, resolved.Bootcamp.ID)
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM coderz.mentee_day_assignments mda
		USING coderz.bootcamp_enrollments be
		WHERE mda.bootcamp_enrollment_id = be.id
		  AND be.bootcamp_id = $1
		  AND mda.weekday = $2
	`, resolved.Bootcamp.ID, normalizedDay); err != nil {
		return nil, err
	}

	for _, username := range targets {
		enrollmentID, ok := enrollmentMap[username]
		if !ok {
			return nil, errors.New("MENTEE_NOT_FOUND")
		}

		if _, err := tx.Exec(ctx, `
			INSERT INTO coderz.mentee_day_assignments (
				bootcamp_enrollment_id,
				weekday,
				created_by
			) VALUES (
				$1,
				$2,
				$3
			)
		`, enrollmentID, normalizedDay, resolved.MemberID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.GetDayAssignments(ctx, userID, normalizedDay)
}

func (s *Service) CreateAssignments(ctx context.Context, userID string, req CreateAssignmentsRequest) (*CreateAssignmentsData, error) {
	resolved, err := s.resolveMentorContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	catalog, ok := findSheet(req.SheetKey)
	if !ok {
		return nil, errors.New("SHEET_NOT_FOUND")
	}

	questionIDs := dedupeStrings(req.QuestionIDs)
	if len(questionIDs) == 0 {
		return nil, errors.New("QUESTION_IDS_REQUIRED")
	}

	selectedQuestions := make([]sheetQuestion, 0, len(questionIDs))
	for _, questionID := range questionIDs {
		question, found := findSheetQuestion(req.SheetKey, questionID)
		if !found {
			return nil, errors.New("QUESTION_NOT_FOUND")
		}
		selectedQuestions = append(selectedQuestions, question)
	}

	targetUsernames := dedupeLower(req.MenteeUsernames)
	if len(targetUsernames) == 0 {
		return nil, errors.New("MENTEES_REQUIRED")
	}

	mentees, err := s.listMenteeRecords(ctx, resolved.Bootcamp.ID)
	if err != nil {
		return nil, err
	}
	menteeByUsername := make(map[string]menteeRecord, len(mentees))
	for _, mentee := range mentees {
		menteeByUsername[strings.ToLower(mentee.Username)] = mentee
	}

	targets := make([]menteeRecord, 0, len(targetUsernames))
	for _, username := range targetUsernames {
		mentee, ok := menteeByUsername[username]
		if !ok {
			return nil, errors.New("MENTEE_NOT_FOUND")
		}
		targets = append(targets, mentee)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	groupTitle := fmt.Sprintf("Algo Buddy %s", catalog.Name)
	if req.Day != "" {
		groupTitle = fmt.Sprintf("Algo Buddy %s %s", capitalizeWord(normalizeDay(req.Day)), catalog.Name)
	}

	description := fmt.Sprintf("App assignment generated from %s", catalog.Name)
	if req.Day != "" {
		description = fmt.Sprintf("App assignment generated for %s from %s", normalizeDay(req.Day), catalog.Name)
	}

	var assignmentGroupID string
	if err := tx.QueryRow(ctx, `
		INSERT INTO coderz.assignment_groups (
			bootcamp_id,
			created_by,
			title,
			description
		) VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING id::text
	`, resolved.Bootcamp.ID, resolved.MemberID, groupTitle, description).Scan(&assignmentGroupID); err != nil {
		return nil, err
	}

	problemIDs := make([]string, 0, len(selectedQuestions))
	for index, question := range selectedQuestions {
		problemID, err := s.getOrCreateProblem(ctx, tx, resolved.Organization.ID, resolved.MemberID, req.SheetKey, question)
		if err != nil {
			return nil, err
		}
		problemIDs = append(problemIDs, problemID)

		if _, err := tx.Exec(ctx, `
			INSERT INTO coderz.assignment_group_problems (
				assignment_group_id,
				problem_id,
				position
			) VALUES (
				$1,
				$2,
				$3
			)
		`, assignmentGroupID, problemID, index+1); err != nil {
			return nil, err
		}
	}

	for _, mentee := range targets {
		var assignmentID string
		if err := tx.QueryRow(ctx, `
			INSERT INTO coderz.assignments (
				assignment_group_id,
				bootcamp_enrollment_id,
				assigned_by,
				status
			) VALUES (
				$1,
				$2,
				$3,
				'active'
			)
			RETURNING id::text
		`, assignmentGroupID, mentee.EnrollmentID, resolved.MemberID).Scan(&assignmentID); err != nil {
			return nil, err
		}

		for _, problemID := range problemIDs {
			if _, err := tx.Exec(ctx, `
				INSERT INTO coderz.assignment_problems (
					assignment_id,
					problem_id,
					status,
					app_progress_status
				) VALUES (
					$1,
					$2,
					'pending',
					'not_started'
				)
			`, assignmentID, problemID); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	if err := s.refreshLeaderboard(ctx, s.pool, resolved.Bootcamp.ID); err != nil {
		return nil, err
	}

	return &CreateAssignmentsData{
		AssignmentGroupID: assignmentGroupID,
		AssignmentsCount:  len(targets),
	}, nil
}

func (s *Service) ListMenteeQuestions(ctx context.Context, userID, username string) ([]QuestionData, error) {
	resolved, err := s.resolveContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	if resolved.Bootcamp == nil || resolved.AccountStatus != "approved" {
		return nil, errors.New("ACCESS_DENIED")
	}

	rows, err := s.pool.Query(ctx, `
		SELECT
			ap.id::text,
			a.id::text,
			COALESCE(u.username, ''),
			p.title,
			COALESCE(p.description, ''),
			p.difficulty::text,
			COALESCE(p.external_link, ''),
			COALESCE(ap.app_progress_status, ''),
			ap.status::text,
			COALESCE(ap.notes, ''),
			COALESCE(ap.resources, ''),
			a.assigned_at,
			ap.completed_at
		FROM coderz.assignment_problems ap
		JOIN coderz.assignments a ON a.id = ap.assignment_id AND a.archived_at IS NULL
		JOIN coderz.problems p ON p.id = ap.problem_id
		JOIN coderz.bootcamp_enrollments be ON be.id = a.bootcamp_enrollment_id
		JOIN coderz.organization_members om ON om.id = be.organization_member_id
		JOIN coderz.users u ON u.id = om.user_id
		WHERE be.bootcamp_id = $1
		  AND be.role = 'mentee'
		  AND be.status = 'active'
		  AND LOWER(COALESCE(u.username, '')) = LOWER($2)
		ORDER BY a.assigned_at DESC, ap.created_at ASC
	`, resolved.Bootcamp.ID, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := make([]QuestionData, 0)
	for rows.Next() {
		row, err := scanQuestionRow(rows)
		if err != nil {
			return nil, err
		}
		questions = append(questions, row.toQuestionData())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

func (s *Service) GetMenteeQuestion(ctx context.Context, userID, username, assignmentProblemID string) (*QuestionData, error) {
	resolved, err := s.resolveContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	if resolved.Bootcamp == nil || resolved.AccountStatus != "approved" {
		return nil, errors.New("ACCESS_DENIED")
	}

	row, err := s.getQuestionRow(ctx, s.pool, resolved.Bootcamp.ID, username, assignmentProblemID)
	if err != nil {
		return nil, err
	}

	data := row.toQuestionData()
	return &data, nil
}

func (s *Service) UpdateMenteeQuestion(ctx context.Context, userID, username, assignmentProblemID string, req UpdateQuestionRequest) (*QuestionData, error) {
	resolved, err := s.resolveContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	if resolved.Bootcamp == nil || resolved.AccountStatus != "approved" {
		return nil, errors.New("ACCESS_DENIED")
	}
	if resolved.Role != "mentor" && !strings.EqualFold(resolved.User.Username, username) {
		return nil, errors.New("ACCESS_DENIED")
	}
	if req.ProgressStatus == nil && req.Solution == nil && req.Resources == nil {
		return nil, errors.New("NO_FIELDS_TO_UPDATE")
	}

	currentRow, err := s.getQuestionRow(ctx, s.pool, resolved.Bootcamp.ID, username, assignmentProblemID)
	if err != nil {
		return nil, err
	}

	progressStatus := currentRow.normalizedProgressStatus()
	if req.ProgressStatus != nil {
		progressStatus = *req.ProgressStatus
	}

	legacyStatus := mapProgressToLegacyStatus(progressStatus)
	setCompletedAt := req.ProgressStatus != nil && progressStatus == "completed"
	clearCompletedAt := req.ProgressStatus != nil && currentRow.normalizedProgressStatus() == "completed" && progressStatus != "completed"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, `
		UPDATE coderz.assignment_problems
		SET
			app_progress_status = $2,
			status = $3,
			notes = CASE WHEN $4 THEN $5 ELSE notes END,
			resources = CASE WHEN $6 THEN $7 ELSE resources END,
			completed_at = CASE
				WHEN $8 THEN CURRENT_TIMESTAMP
				WHEN $9 THEN NULL
				ELSE completed_at
			END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, assignmentProblemID, progressStatus, legacyStatus, req.Solution != nil, valueOrEmpty(req.Solution), req.Resources != nil, valueOrEmpty(req.Resources), setCompletedAt, clearCompletedAt); err != nil {
		return nil, err
	}

	if err := s.updateAssignmentAggregate(ctx, tx, currentRow.AssignmentID); err != nil {
		return nil, err
	}
	if err := s.refreshLeaderboard(ctx, tx, resolved.Bootcamp.ID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	updatedRow, err := s.getQuestionRow(ctx, s.pool, resolved.Bootcamp.ID, username, assignmentProblemID)
	if err != nil {
		return nil, err
	}

	data := updatedRow.toQuestionData()
	return &data, nil
}

func (s *Service) GetMenteeProfile(ctx context.Context, userID, username string) (*ProfileData, error) {
	resolved, err := s.resolveContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	if resolved.Bootcamp == nil || resolved.AccountStatus != "approved" {
		return nil, errors.New("ACCESS_DENIED")
	}

	mentee, err := s.findMenteeByUsername(ctx, s.pool, resolved.Bootcamp.ID, username)
	if err != nil {
		return nil, err
	}

	solved, err := s.countCompletedProblems(ctx, resolved.Bootcamp.ID, username)
	if err != nil {
		return nil, err
	}

	var (
		bio      string
		github   string
		linkedin string
	)
	if err := s.pool.QueryRow(ctx, `
		SELECT
			COALESCE(bio, ''),
			COALESCE(github_url, ''),
			COALESCE(linkedin_url, '')
		FROM coderz.users
		WHERE id = $1
	`, mentee.UserID).Scan(&bio, &github, &linkedin); err != nil {
		return nil, err
	}

	return &ProfileData{
		FirstName: mentee.FirstName,
		LastName:  mentee.LastName,
		Username:  mentee.Username,
		Email:     "",
		Solved:    solved,
		JoinedAt:  mentee.EnrolledAt.Format(time.RFC3339),
		Bio:       bio,
		Github:    github,
		Linkedin:  linkedin,
	}, nil
}

func (s *Service) GetMyProfile(ctx context.Context, userID string) (*ProfileData, error) {
	resolved, err := s.resolveContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	solved := 0
	if resolved.Bootcamp != nil {
		solved, err = s.countCompletedProblems(ctx, resolved.Bootcamp.ID, resolved.User.Username)
		if err != nil {
			return nil, err
		}
	}

	var createdAt time.Time
	if err := s.pool.QueryRow(ctx, `
		SELECT created_at
		FROM coderz.users
		WHERE id = $1
	`, resolved.UserID).Scan(&createdAt); err != nil {
		return nil, err
	}
	joinedAt := createdAt.Format(time.RFC3339)

	if resolved.EnrollmentID != "" {
		var enrolledAt time.Time
		if err := s.pool.QueryRow(ctx, `
			SELECT enrolled_at
			FROM coderz.bootcamp_enrollments
			WHERE id = $1
		`, resolved.EnrollmentID).Scan(&enrolledAt); err == nil {
			joinedAt = enrolledAt.Format(time.RFC3339)
		}
	}

	return &ProfileData{
		FirstName: resolved.User.FirstName,
		LastName:  resolved.User.LastName,
		Username:  resolved.User.Username,
		Email:     resolved.User.Email,
		Solved:    solved,
		JoinedAt:  joinedAt,
		Bio:       resolved.User.Bio,
		Github:    resolved.User.Github,
		Linkedin:  resolved.User.Linkedin,
	}, nil
}

func (s *Service) UpdateMyProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*ProfileData, error) {
	resolved, err := s.resolveContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	username := normalizeUsername(req.Username)
	if !usernamePattern.MatchString(username) {
		return nil, errors.New("INVALID_USERNAME")
	}

	fullName := strings.TrimSpace(strings.TrimSpace(req.FirstName) + " " + strings.TrimSpace(req.LastName))
	if _, err := s.pool.Exec(ctx, `
		UPDATE coderz.users
		SET
			name = $2,
			email = $3,
			username = $4,
			bio = NULLIF($5, ''),
			github_url = NULLIF($6, ''),
			linkedin_url = NULLIF($7, ''),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, resolved.UserID, fullName, req.Email, username, req.Bio, req.Github, req.Linkedin); err != nil {
		lowerErr := strings.ToLower(err.Error())
		if strings.Contains(lowerErr, "uq_users_username") {
			return nil, errors.New("USERNAME_ALREADY_EXISTS")
		}
		if strings.Contains(lowerErr, "users_email_key") {
			return nil, errors.New("EMAIL_ALREADY_EXISTS")
		}
		return nil, err
	}

	return s.GetMyProfile(ctx, userID)
}

func (s *Service) UpdateMyPassword(ctx context.Context, userID string, req UpdatePasswordRequest) error {
	if !validatePasswordComplexity(req.NewPassword) {
		return errors.New("PASSWORD_MUST_CONTAIN_LETTER_AND_NUMBER")
	}

	var passwordHash pgtype.Text
	if err := s.pool.QueryRow(ctx, `
		SELECT password_hash
		FROM coderz.users
		WHERE id = $1
	`, userID).Scan(&passwordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("USER_NOT_FOUND")
		}
		return err
	}
	if !passwordHash.Valid {
		return errors.New("PASSWORD_LOGIN_NOT_AVAILABLE")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash.String), []byte(req.CurrentPassword)); err != nil {
		return errors.New("INVALID_CURRENT_PASSWORD")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, `
		UPDATE coderz.users
		SET
			password_hash = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, userID, string(newHash)); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		DELETE FROM coderz.refresh_tokens
		WHERE user_id = $1
	`, userID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) GetLeaderboard(ctx context.Context, userID string) ([]LeaderboardEntryData, error) {
	resolved, err := s.resolveContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	if resolved.Bootcamp == nil || resolved.AccountStatus != "approved" {
		return nil, errors.New("ACCESS_DENIED")
	}

	if err := s.refreshLeaderboard(ctx, s.pool, resolved.Bootcamp.ID); err != nil {
		return nil, err
	}

	rows, err := s.pool.Query(ctx, `
		SELECT
			COALESCE(u.username, ''),
			u.name,
			COALESCE(le.problems_completed, 0) AS solved
		FROM coderz.bootcamp_enrollments be
		JOIN coderz.organization_members om ON om.id = be.organization_member_id
		JOIN coderz.users u ON u.id = om.user_id
		LEFT JOIN coderz.leaderboard_entries le
			ON le.bootcamp_enrollment_id = be.id
			AND le.bootcamp_id = be.bootcamp_id
		WHERE be.bootcamp_id = $1
		  AND be.role = 'mentee'
		  AND be.status = 'active'
		ORDER BY COALESCE(le.rank, 2147483647), solved DESC, u.name ASC
	`, resolved.Bootcamp.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]LeaderboardEntryData, 0)
	for rows.Next() {
		var (
			username string
			fullName string
			solved   int
		)
		if err := rows.Scan(&username, &fullName, &solved); err != nil {
			return nil, err
		}

		firstName, lastName := splitName(fullName)
		entries = append(entries, LeaderboardEntryData{
			Username:  username,
			FirstName: firstName,
			LastName:  lastName,
			Solved:    solved,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *Service) resolveContext(ctx context.Context, userID string) (*resolvedContext, error) {
	var (
		foundID  string
		name     string
		username string
		email    string
		bio      string
		github   string
		linkedin string
	)
	if err := s.pool.QueryRow(ctx, `
		SELECT
			id::text,
			name,
			COALESCE(username, ''),
			COALESCE(email, ''),
			COALESCE(bio, ''),
			COALESCE(github_url, ''),
			COALESCE(linkedin_url, '')
		FROM coderz.users
		WHERE id = $1
	`, userID).Scan(&foundID, &name, &username, &email, &bio, &github, &linkedin); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("USER_NOT_FOUND")
		}
		return nil, err
	}

	firstName, lastName := splitName(name)
	resolved := &resolvedContext{
		UserID: foundID,
		User: UserData{
			ID:        foundID,
			Name:      name,
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
			Email:     email,
			Bio:       bio,
			Github:    github,
			Linkedin:  linkedin,
		},
		Role:          "unknown",
		AccountStatus: "unassigned",
	}

	var (
		memberID       string
		orgRole        string
		orgID          string
		orgName        string
		orgSlug        string
		bootcampID     string
		bootcampName   string
		enrollmentID   string
		enrollmentRole string
		assignedSheet  string
	)
	err := s.pool.QueryRow(ctx, `
		SELECT
			om.id::text,
			om.role::text,
			o.id::text,
			o.name,
			o.slug,
			b.id::text,
			b.name,
			be.id::text,
			be.role::text,
			COALESCE(be.assigned_sheet_key, '')
		FROM coderz.bootcamp_enrollments be
		JOIN coderz.organization_members om ON om.id = be.organization_member_id
		JOIN coderz.organizations o ON o.id = om.organization_id
		JOIN coderz.bootcamps b ON b.id = be.bootcamp_id
		WHERE om.user_id = $1
		  AND o.status = 'approved'
		  AND b.archived_at IS NULL
		  AND b.is_active = TRUE
		  AND be.status = 'active'
		ORDER BY b.created_at DESC, be.enrolled_at DESC
		LIMIT 1
	`, userID).Scan(&memberID, &orgRole, &orgID, &orgName, &orgSlug, &bootcampID, &bootcampName, &enrollmentID, &enrollmentRole, &assignedSheet)
	if err == nil {
		resolved.MemberID = memberID
		resolved.OrgRole = orgRole
		resolved.Organization = &OrganizationData{ID: orgID, Name: orgName, Slug: orgSlug}
		resolved.Bootcamp = &BootcampData{ID: bootcampID, Name: bootcampName}
		resolved.EnrollmentID = enrollmentID
		resolved.EnrollmentRole = enrollmentRole
		resolved.AssignedSheet = assignedSheet
		resolved.AccountStatus = "approved"
		if enrollmentRole == "mentor" || orgRole == "admin" || orgRole == "mentor" {
			resolved.Role = "mentor"
		} else {
			resolved.Role = "mentee"
		}
		return resolved, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	err = s.pool.QueryRow(ctx, `
		SELECT
			om.id::text,
			om.role::text,
			o.id::text,
			o.name,
			o.slug,
			b.id::text,
			b.name
		FROM coderz.organization_members om
		JOIN coderz.organizations o ON o.id = om.organization_id
		JOIN coderz.bootcamps b ON b.organization_id = o.id
		WHERE om.user_id = $1
		  AND o.status = 'approved'
		  AND om.role IN ('admin', 'mentor')
		  AND b.archived_at IS NULL
		  AND b.is_active = TRUE
		ORDER BY b.created_at DESC, om.joined_at DESC
		LIMIT 1
	`, userID).Scan(&memberID, &orgRole, &orgID, &orgName, &orgSlug, &bootcampID, &bootcampName)
	if err == nil {
		resolved.MemberID = memberID
		resolved.OrgRole = orgRole
		resolved.Organization = &OrganizationData{ID: orgID, Name: orgName, Slug: orgSlug}
		resolved.Bootcamp = &BootcampData{ID: bootcampID, Name: bootcampName}
		resolved.Role = "mentor"
		resolved.AccountStatus = "approved"
		return resolved, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	var status string
	err = s.pool.QueryRow(ctx, `
		SELECT
			mr.status,
			COALESCE(mr.sheet_key, ''),
			o.id::text,
			o.name,
			o.slug,
			b.id::text,
			b.name
		FROM coderz.mentee_requests mr
		JOIN coderz.organizations o ON o.id = mr.organization_id
		JOIN coderz.bootcamps b ON b.id = mr.bootcamp_id
		WHERE mr.user_id = $1
		ORDER BY mr.created_at DESC
		LIMIT 1
	`, userID).Scan(&status, &assignedSheet, &orgID, &orgName, &orgSlug, &bootcampID, &bootcampName)
	if err == nil {
		resolved.Organization = &OrganizationData{ID: orgID, Name: orgName, Slug: orgSlug}
		resolved.Bootcamp = &BootcampData{ID: bootcampID, Name: bootcampName}
		resolved.Role = "mentee"
		resolved.AccountStatus = status
		resolved.AssignedSheet = assignedSheet
		return resolved, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return resolved, nil
}

func (s *Service) resolveMentorContext(ctx context.Context, userID string) (*resolvedContext, error) {
	resolved, err := s.resolveContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	if resolved.Role != "mentor" || resolved.AccountStatus != "approved" || resolved.Organization == nil || resolved.Bootcamp == nil || resolved.MemberID == "" {
		return nil, errors.New("ACCESS_DENIED")
	}
	return resolved, nil
}

func (s *Service) getDefaultSignupContext(ctx context.Context) (*OrganizationData, *BootcampData, error) {
	var (
		orgID        string
		orgName      string
		orgSlug      string
		bootcampID   string
		bootcampName string
	)
	if err := s.pool.QueryRow(ctx, `
		SELECT
			o.id::text,
			o.name,
			o.slug,
			b.id::text,
			b.name
		FROM coderz.bootcamps b
		JOIN coderz.organizations o ON o.id = b.organization_id
		WHERE o.status = 'approved'
		  AND b.archived_at IS NULL
		  AND b.is_active = TRUE
		ORDER BY b.created_at DESC
		LIMIT 1
	`).Scan(&orgID, &orgName, &orgSlug, &bootcampID, &bootcampName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, errors.New("BOOTCAMP_NOT_CONFIGURED")
		}
		return nil, nil, err
	}

	return &OrganizationData{ID: orgID, Name: orgName, Slug: orgSlug}, &BootcampData{ID: bootcampID, Name: bootcampName}, nil
}

func (s *Service) ensureOrganizationMember(ctx context.Context, q db.DBTX, organizationID, userID string) (string, error) {
	var memberID string
	err := q.QueryRow(ctx, `
		SELECT id::text
		FROM coderz.organization_members
		WHERE organization_id = $1
		  AND user_id = $2
		LIMIT 1
	`, organizationID, userID).Scan(&memberID)
	if err == nil {
		return memberID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}

	if err := q.QueryRow(ctx, `
		INSERT INTO coderz.organization_members (
			organization_id,
			user_id,
			role
		) VALUES (
			$1,
			$2,
			'mentee'
		)
		RETURNING id::text
	`, organizationID, userID).Scan(&memberID); err != nil {
		return "", err
	}

	return memberID, nil
}

func (s *Service) ensureBootcampEnrollment(ctx context.Context, q db.DBTX, bootcampID, memberID, assignedSheet string) error {
	var enrollmentID string
	err := q.QueryRow(ctx, `
		SELECT id::text
		FROM coderz.bootcamp_enrollments
		WHERE bootcamp_id = $1
		  AND organization_member_id = $2
		LIMIT 1
	`, bootcampID, memberID).Scan(&enrollmentID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		_, err = q.Exec(ctx, `
			INSERT INTO coderz.bootcamp_enrollments (
				bootcamp_id,
				organization_member_id,
				role,
				status,
				assigned_sheet_key
			) VALUES (
				$1,
				$2,
				'mentee',
				'active',
				NULLIF($3, '')
			)
		`, bootcampID, memberID, assignedSheet)
		return err
	}

	_, err = q.Exec(ctx, `
		UPDATE coderz.bootcamp_enrollments
		SET
			role = 'mentee',
			status = 'active',
			assigned_sheet_key = NULLIF($2, '')
		WHERE id = $1
	`, enrollmentID, assignedSheet)
	return err
}

func (s *Service) listMenteeEnrollmentMap(ctx context.Context, q db.DBTX, bootcampID string) (map[string]string, error) {
	rows, err := q.Query(ctx, `
		SELECT
			LOWER(COALESCE(u.username, '')) AS username,
			be.id::text
		FROM coderz.bootcamp_enrollments be
		JOIN coderz.organization_members om ON om.id = be.organization_member_id
		JOIN coderz.users u ON u.id = om.user_id
		WHERE be.bootcamp_id = $1
		  AND be.role = 'mentee'
		  AND be.status = 'active'
	`, bootcampID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var username string
		var enrollmentID string
		if err := rows.Scan(&username, &enrollmentID); err != nil {
			return nil, err
		}
		result[username] = enrollmentID
	}

	return result, rows.Err()
}

func (s *Service) listMenteeRecords(ctx context.Context, bootcampID string) ([]menteeRecord, error) {
	return s.listMenteeRecordsWithQuery(ctx, s.pool, bootcampID)
}

func (s *Service) listMenteeRecordsWithQuery(ctx context.Context, q db.DBTX, bootcampID string) ([]menteeRecord, error) {
	rows, err := q.Query(ctx, `
		SELECT
			be.id::text,
			om.id::text,
			u.id::text,
			u.name,
			COALESCE(u.username, ''),
			COALESCE(u.email, ''),
			COALESCE(be.assigned_sheet_key, ''),
			be.enrolled_at
		FROM coderz.bootcamp_enrollments be
		JOIN coderz.organization_members om ON om.id = be.organization_member_id
		JOIN coderz.users u ON u.id = om.user_id
		WHERE be.bootcamp_id = $1
		  AND be.role = 'mentee'
		  AND be.status = 'active'
		ORDER BY u.name ASC
	`, bootcampID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mentees := make([]menteeRecord, 0)
	for rows.Next() {
		var (
			enrollmentID  string
			memberID      string
			userID        string
			fullName      string
			username      string
			email         string
			assignedSheet string
			enrolledAt    time.Time
		)
		if err := rows.Scan(&enrollmentID, &memberID, &userID, &fullName, &username, &email, &assignedSheet, &enrolledAt); err != nil {
			return nil, err
		}

		firstName, lastName := splitName(fullName)
		mentees = append(mentees, menteeRecord{
			EnrollmentID:  enrollmentID,
			MemberID:      memberID,
			UserID:        userID,
			FirstName:     firstName,
			LastName:      lastName,
			Username:      username,
			Email:         email,
			AssignedSheet: assignedSheet,
			EnrolledAt:    enrolledAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mentees, nil
}

func (s *Service) findMenteeByUsername(ctx context.Context, q db.DBTX, bootcampID, username string) (*menteeRecord, error) {
	var (
		enrollmentID  string
		memberID      string
		userID        string
		fullName      string
		foundUsername string
		email         string
		assignedSheet string
		enrolledAt    time.Time
	)
	if err := q.QueryRow(ctx, `
		SELECT
			be.id::text,
			om.id::text,
			u.id::text,
			u.name,
			COALESCE(u.username, ''),
			COALESCE(u.email, ''),
			COALESCE(be.assigned_sheet_key, ''),
			be.enrolled_at
		FROM coderz.bootcamp_enrollments be
		JOIN coderz.organization_members om ON om.id = be.organization_member_id
		JOIN coderz.users u ON u.id = om.user_id
		WHERE be.bootcamp_id = $1
		  AND be.role = 'mentee'
		  AND be.status = 'active'
		  AND LOWER(COALESCE(u.username, '')) = LOWER($2)
		LIMIT 1
	`, bootcampID, username).Scan(&enrollmentID, &memberID, &userID, &fullName, &foundUsername, &email, &assignedSheet, &enrolledAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("MENTEE_NOT_FOUND")
		}
		return nil, err
	}

	firstName, lastName := splitName(fullName)
	return &menteeRecord{
		EnrollmentID:  enrollmentID,
		MemberID:      memberID,
		UserID:        userID,
		FirstName:     firstName,
		LastName:      lastName,
		Username:      foundUsername,
		Email:         email,
		AssignedSheet: assignedSheet,
		EnrolledAt:    enrolledAt,
	}, nil
}

func (s *Service) getOrCreateProblem(ctx context.Context, q db.DBTX, organizationID, createdBy, sheetKey string, question sheetQuestion) (string, error) {
	link := catalogLink(sheetKey, question.ID)
	var problemID string
	err := q.QueryRow(ctx, `
		SELECT id::text
		FROM coderz.problems
		WHERE organization_id = $1
		  AND external_link = $2
		LIMIT 1
	`, organizationID, link).Scan(&problemID)
	if err == nil {
		return problemID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}

	if err := q.QueryRow(ctx, `
		INSERT INTO coderz.problems (
			organization_id,
			created_by,
			title,
			description,
			difficulty,
			external_link
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
		RETURNING id::text
	`, organizationID, createdBy, question.Title, question.Description, question.Difficulty, link).Scan(&problemID); err != nil {
		return "", err
	}

	return problemID, nil
}

func (s *Service) getQuestionRow(ctx context.Context, q db.DBTX, bootcampID, username, assignmentProblemID string) (*questionRow, error) {
	row := q.QueryRow(ctx, `
		SELECT
			ap.id::text,
			a.id::text,
			COALESCE(u.username, ''),
			p.title,
			COALESCE(p.description, ''),
			p.difficulty::text,
			COALESCE(p.external_link, ''),
			COALESCE(ap.app_progress_status, ''),
			ap.status::text,
			COALESCE(ap.notes, ''),
			COALESCE(ap.resources, ''),
			a.assigned_at,
			ap.completed_at
		FROM coderz.assignment_problems ap
		JOIN coderz.assignments a ON a.id = ap.assignment_id AND a.archived_at IS NULL
		JOIN coderz.problems p ON p.id = ap.problem_id
		JOIN coderz.bootcamp_enrollments be ON be.id = a.bootcamp_enrollment_id
		JOIN coderz.organization_members om ON om.id = be.organization_member_id
		JOIN coderz.users u ON u.id = om.user_id
		WHERE ap.id = $1
		  AND be.bootcamp_id = $2
		  AND LOWER(COALESCE(u.username, '')) = LOWER($3)
		LIMIT 1
	`, assignmentProblemID, bootcampID, username)

	question, err := scanQuestionRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("QUESTION_NOT_FOUND")
		}
		return nil, err
	}

	return &question, nil
}

func (s *Service) updateAssignmentAggregate(ctx context.Context, q db.DBTX, assignmentID string) error {
	var totalCount int
	var completedCount int
	if err := q.QueryRow(ctx, `
		SELECT
			COUNT(*)::int,
			COUNT(*) FILTER (WHERE status = 'completed')::int
		FROM coderz.assignment_problems
		WHERE assignment_id = $1
	`, assignmentID).Scan(&totalCount, &completedCount); err != nil {
		return err
	}

	status := "active"
	if totalCount > 0 && totalCount == completedCount {
		status = "completed"
	}

	_, err := q.Exec(ctx, `
		UPDATE coderz.assignments
		SET
			status = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, assignmentID, status)
	return err
}

func (s *Service) refreshLeaderboard(ctx context.Context, q db.DBTX, bootcampID string) error {
	rows, err := q.Query(ctx, `
		SELECT
			be.id::text,
			COUNT(ap.id)::int AS total_assigned,
			COUNT(*) FILTER (
				WHERE ap.app_progress_status = 'completed'
				   OR ap.status = 'completed'
			)::int AS completed_count,
			COUNT(*) FILTER (
				WHERE ap.app_progress_status <> 'not_started'
				   OR ap.status IN ('attempted', 'completed')
			)::int AS attempted_count
		FROM coderz.bootcamp_enrollments be
		JOIN coderz.organization_members om ON om.id = be.organization_member_id
		JOIN coderz.users u ON u.id = om.user_id
		LEFT JOIN coderz.assignments a
			ON a.bootcamp_enrollment_id = be.id
			AND a.archived_at IS NULL
		LEFT JOIN coderz.assignment_problems ap ON ap.assignment_id = a.id
		WHERE be.bootcamp_id = $1
		  AND be.role = 'mentee'
		  AND be.status = 'active'
		GROUP BY be.id, u.username, u.name
		ORDER BY completed_count DESC, attempted_count DESC, u.name ASC
	`, bootcampID)
	if err != nil {
		return err
	}
	defer rows.Close()

	type leaderboardRow struct {
		enrollmentID string
		total        int
		completed    int
		attempted    int
	}
	stats := make([]leaderboardRow, 0)
	for rows.Next() {
		var item leaderboardRow
		if err := rows.Scan(&item.enrollmentID, &item.total, &item.completed, &item.attempted); err != nil {
			return err
		}
		stats = append(stats, item)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if _, err := q.Exec(ctx, `
		DELETE FROM coderz.leaderboard_entries
		WHERE bootcamp_id = $1
	`, bootcampID); err != nil {
		return err
	}

	for index, item := range stats {
		completionRate := float32(0)
		if item.total > 0 {
			completionRate = float32(item.completed) / float32(item.total)
		}
		score := (item.completed * 10) + (item.attempted * 3)

		if _, err := q.Exec(ctx, `
			INSERT INTO coderz.leaderboard_entries (
				bootcamp_id,
				bootcamp_enrollment_id,
				problems_completed,
				problems_attempted,
				completion_rate,
				streak_days,
				score,
				rank,
				calculated_at
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				0,
				$6,
				$7,
				CURRENT_TIMESTAMP
			)
		`, bootcampID, item.enrollmentID, item.completed, item.attempted, completionRate, score, index+1); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) countCompletedProblems(ctx context.Context, bootcampID, username string) (int, error) {
	var solved int
	if err := s.pool.QueryRow(ctx, `
		SELECT COUNT(*)::int
		FROM coderz.assignment_problems ap
		JOIN coderz.assignments a ON a.id = ap.assignment_id AND a.archived_at IS NULL
		JOIN coderz.bootcamp_enrollments be ON be.id = a.bootcamp_enrollment_id
		JOIN coderz.organization_members om ON om.id = be.organization_member_id
		JOIN coderz.users u ON u.id = om.user_id
		WHERE be.bootcamp_id = $1
		  AND LOWER(COALESCE(u.username, '')) = LOWER($2)
		  AND (
			ap.app_progress_status = 'completed'
			OR ap.status = 'completed'
		  )
	`, bootcampID, username).Scan(&solved); err != nil {
		return 0, err
	}
	return solved, nil
}

func scanQuestionRow(scanner interface{ Scan(dest ...any) error }) (questionRow, error) {
	var row questionRow
	err := scanner.Scan(
		&row.ID,
		&row.AssignmentID,
		&row.TargetUsername,
		&row.Title,
		&row.Description,
		&row.Difficulty,
		&row.ExternalLink,
		&row.AppProgress,
		&row.LegacyStatus,
		&row.Notes,
		&row.Resources,
		&row.AssignedAt,
		&row.CompletedAt,
	)
	return row, err
}

func (q questionRow) normalizedProgressStatus() string {
	if q.AppProgress != "" {
		return q.AppProgress
	}
	switch q.LegacyStatus {
	case "completed":
		return "completed"
	case "attempted":
		return "revision_needed"
	default:
		return "not_started"
	}
}

func (q questionRow) toQuestionData() QuestionData {
	description := q.Description
	topic := "General"
	if catalogQuestion, ok := findSheetQuestionByLink(q.ExternalLink); ok {
		description = catalogQuestion.Description
		topic = catalogQuestion.Topic
	}

	progressStatus := q.normalizedProgressStatus()
	status := "pending"
	if progressStatus == "completed" {
		status = "completed"
	}

	completedAt := ""
	if q.CompletedAt.Valid {
		completedAt = q.CompletedAt.Time.Format(time.RFC3339)
	}

	return QuestionData{
		ID:             q.ID,
		Title:          q.Title,
		Description:    description,
		Difficulty:     q.Difficulty,
		Topic:          topic,
		Status:         status,
		ProgressStatus: progressStatus,
		AssignedAt:     q.AssignedAt.Format(time.RFC3339),
		CompletedAt:    completedAt,
		Solution:       q.Notes,
		Resources:      q.Resources,
	}
}

func splitName(name string) (string, string) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", ""
	}
	parts := strings.Fields(trimmed)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func validatePasswordComplexity(password string) bool {
	hasLetter := false
	hasNumber := false

	for _, char := range password {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
		}
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
		if hasLetter && hasNumber {
			return true
		}
	}

	return false
}

func dedupeLower(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result
}

func dedupeStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		normalized := strings.TrimSpace(value)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result
}

func mapProgressToLegacyStatus(progress string) string {
	switch progress {
	case "completed":
		return "completed"
	case "discussion_needed", "revision_needed":
		return "attempted"
	default:
		return "pending"
	}
}

func normalizeDay(day string) string {
	return strings.ToLower(strings.TrimSpace(day))
}

func capitalizeWord(value string) string {
	if value == "" {
		return ""
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
