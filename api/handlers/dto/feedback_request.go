package dto

import "time"

type CreateQuestionRequest struct {
	Question string `json:"question" example:"Apa pendapat Anda tentang pelatihan ini?"`
}

type SubmitAnswerRequest struct {
	QuestionID string `json:"question_id" example:"b5a1c6c3-1234-4bcd-9123-a12b34cd56ef"`
	Answer     string `json:"answer" example:"Sangat bermanfaat dan jelas"`
}


type FeedbackStudentResponse struct {
	ID        string    `json:"id" example:"aa5bada7-1063-4817-b31d-3a62f233e20f"`
	Name      string    `json:"name" example:"Guru 2"`
	Email     string    `json:"email" example:"guru2@gmail.com"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2025-10-31T10:06:20.249632+07:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-10-31T16:43:55.96876+07:00"`
}

type FeedbackAnswerWithStudentResponse struct {
	ID         string                 `json:"id" example:"abf397fe-76b9-4222-8d1e-2d51c6be6f9b"`
	QuestionID string                 `json:"question_id" example:"0bd98683-6f80-47d8-9017-b15106ba9b53"`
	StudentID  string                 `json:"student_id" example:"aa5bada7-1063-4817-b31d-3a62f233e20f"`
	Answer     string                 `json:"answer" example:"halo mas"`
	CreatedAt  time.Time              `json:"created_at" example:"2025-10-31T16:28:34.496183+07:00"`
	Student    FeedbackStudentResponse `json:"student"`
}

type FeedbackQuestionWithAnswersResponse struct {
	ID        string                           `json:"id" example:"b5451826-53d0-4904-93e8-3a88e08952f7"`
	Question  string                           `json:"question" example:"Bagaimana kelas baru nyaaaaa?"`
	CreatedBy string                           `json:"created_by" example:"b7dfe843-4297-4d13-b666-d865df01ecbc"`
	CreatedAt time.Time                        `json:"created_at" example:"2025-10-30T17:05:09.051049+07:00"`
	UpdatedAt time.Time                        `json:"updated_at" example:"2025-10-30T17:05:09.051049+07:00"`
	Answers   []FeedbackAnswerWithStudentResponse `json:"answers,omitempty"`
}