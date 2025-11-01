package feedback

import (
	"api-shiners/pkg/entities"
	"context"

	"github.com/google/uuid"
)

type FeedbackService interface {
	CreateQuestion(ctx context.Context, questionText string, createdBy uuid.UUID) (*entities.FeedbackQuestion, error)
	GetAllQuestions(ctx context.Context) ([]entities.FeedbackQuestion, error)
	SubmitAnswer(ctx context.Context, questionID uuid.UUID, studentID uuid.UUID, answer string) error
	// GetStudentAnswers(ctx context.Context, studentID uuid.UUID) ([]entities.FeedbackAnswer, error)
	GetQuestionsWithAnswersByTeacher(ctx context.Context, teacherID uuid.UUID) ([]entities.FeedbackQuestion, error)
}

type feedbackService struct {
	repo FeedbackRepository
}

func NewFeedbackService(repo FeedbackRepository) FeedbackService {
	return &feedbackService{repo}
}

func (s *feedbackService) CreateQuestion(ctx context.Context, questionText string, createdBy uuid.UUID) (*entities.FeedbackQuestion, error) {
	q := &entities.FeedbackQuestion{
		Question:  questionText,
		CreatedBy: createdBy,
	}
	if err := s.repo.CreateQuestion(ctx, q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *feedbackService) GetAllQuestions(ctx context.Context) ([]entities.FeedbackQuestion, error) {
	return s.repo.GetAllQuestions(ctx)
}

func (s *feedbackService) SubmitAnswer(ctx context.Context, questionID uuid.UUID, studentID uuid.UUID, answer string) error {
	a := &entities.FeedbackAnswer{
		QuestionID: questionID,
		StudentID:  studentID,
		Answer:     answer,
	}
	return s.repo.SubmitAnswer(ctx, a)
}


func (s *feedbackService) GetQuestionsWithAnswersByTeacher(ctx context.Context, teacherID uuid.UUID) ([]entities.FeedbackQuestion, error) {
	return s.repo.GetQuestionsWithAnswersByTeacher(ctx, teacherID)
}