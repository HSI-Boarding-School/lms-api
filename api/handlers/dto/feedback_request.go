package dto

type CreateQuestionRequest struct {
	Question string `json:"question" example:"Apa pendapat Anda tentang pelatihan ini?"`
}

type SubmitAnswerRequest struct {
	QuestionID string `json:"question_id" example:"b5a1c6c3-1234-4bcd-9123-a12b34cd56ef"`
	Answer     string `json:"answer" example:"Sangat bermanfaat dan jelas"`
}