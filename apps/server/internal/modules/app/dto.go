package app

type UserData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio,omitempty"`
	Github    string `json:"github,omitempty"`
	Linkedin  string `json:"linkedin,omitempty"`
}

type OrganizationData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type BootcampData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EnrollmentData struct {
	ID            string `json:"id,omitempty"`
	AssignedSheet string `json:"assignedSheet,omitempty"`
}

type ContextData struct {
	Role          string            `json:"role"`
	AccountStatus string            `json:"accountStatus"`
	User          UserData          `json:"user"`
	Organization  *OrganizationData `json:"organization,omitempty"`
	Bootcamp      *BootcampData     `json:"bootcamp,omitempty"`
	Enrollment    *EnrollmentData   `json:"enrollment,omitempty"`
}

type MenteeSignupRequest struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string `json:"lastName" validate:"omitempty,max=50"`
	Username  string `json:"username" validate:"required,min=3,max=80"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=50,password_complexity"`
}

type MenteeSignupData struct {
	RequestID string `json:"requestId"`
	Status    string `json:"status"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type MenteeRequestData struct {
	ID            string `json:"id"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	SignedUpAt    string `json:"signedUpAt"`
	Status        string `json:"status"`
	AssignedSheet string `json:"assignedSheet,omitempty"`
}

type ReviewMenteeRequest struct {
	Status   string `json:"status" validate:"required,oneof=approved rejected"`
	SheetKey string `json:"sheetKey" validate:"omitempty,oneof=gfg-dsa-360 strivers-dsa-sheet"`
}

type SheetQuestionData struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Topic      string `json:"topic"`
	Difficulty string `json:"difficulty"`
}

type SheetData struct {
	Key       string              `json:"key"`
	Name      string              `json:"name"`
	Questions []SheetQuestionData `json:"questions"`
}

type DayAssignmentMenteeData struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Assigned      bool   `json:"assigned"`
	AssignedSheet string `json:"assignedSheet,omitempty"`
}

type DayAssignmentsData struct {
	Day     string                    `json:"day"`
	Mentees []DayAssignmentMenteeData `json:"mentees"`
}

type UpdateDayAssignmentsRequest struct {
	Usernames []string `json:"usernames" validate:"required,min=0,dive,min=3,max=80"`
}

type CreateAssignmentsRequest struct {
	Day             string   `json:"day" validate:"omitempty,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	MenteeUsernames []string `json:"menteeUsernames" validate:"required,min=1,dive,min=3,max=80"`
	SheetKey        string   `json:"sheetKey" validate:"required,oneof=gfg-dsa-360 strivers-dsa-sheet"`
	QuestionIDs     []string `json:"questionIds" validate:"required,min=1,dive,min=1,max=64"`
}

type CreateAssignmentsData struct {
	AssignmentGroupID string `json:"assignmentGroupId"`
	AssignmentsCount  int    `json:"assignmentsCount"`
}

type QuestionData struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Difficulty     string `json:"difficulty"`
	Topic          string `json:"topic"`
	Status         string `json:"status"`
	ProgressStatus string `json:"progressStatus"`
	AssignedAt     string `json:"assignedAt"`
	CompletedAt    string `json:"completedAt,omitempty"`
	Solution       string `json:"solution,omitempty"`
	Resources      string `json:"resources,omitempty"`
}

type UpdateQuestionRequest struct {
	ProgressStatus *string `json:"progressStatus" validate:"omitempty,oneof=not_started discussion_needed revision_needed completed"`
	Solution       *string `json:"solution" validate:"omitempty,max=2000"`
	Resources      *string `json:"resources" validate:"omitempty,max=2000"`
}

type ProfileData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Solved    int    `json:"solved"`
	JoinedAt  string `json:"joinedAt"`
	Bio       string `json:"bio,omitempty"`
	Github    string `json:"github,omitempty"`
	Linkedin  string `json:"linkedin,omitempty"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string `json:"lastName" validate:"omitempty,max=50"`
	Username  string `json:"username" validate:"required,min=3,max=80"`
	Email     string `json:"email" validate:"required,email"`
	Bio       string `json:"bio" validate:"omitempty,max=500"`
	Github    string `json:"github" validate:"omitempty,url"`
	Linkedin  string `json:"linkedin" validate:"omitempty,url"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=8,max=50"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=50,password_complexity"`
}

type LeaderboardEntryData struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Solved    int    `json:"solved"`
}
