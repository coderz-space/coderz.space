package tasks

type AssignQuestionRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Difficulty  string `json:"difficulty" validate:"required"`
	Topic       string `json:"topic"`
}

type UpdateProgressRequest struct {
	ProgressStatus string `json:"progressStatus" validate:"required"`
}

type UpdateDetailsRequest struct {
	Solution  string `json:"solution"`
	Resources string `json:"resources"`
}

type QuestionDTO struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Difficulty     string `json:"difficulty"`
	Topic          string `json:"topic"`
	Status         string `json:"status"`
	ProgressStatus string `json:"progressStatus"`
	AssignedAt     string `json:"assignedAt"`
	CompletedAt    string `json:"completedAt,omitempty"`
	SolutionUrl    string `json:"solutionUrl,omitempty"`
	Solution       string `json:"solution,omitempty"`
	Resources      string `json:"resources,omitempty"`
}
