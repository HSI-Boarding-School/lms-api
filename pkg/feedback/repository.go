package feedback

import (
	"api-shiners/pkg/entities"
	"context"

	"gorm.io/gorm/clause"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeedbackRepository interface {
	CreateQuestion(ctx context.Context, question *entities.FeedbackQuestion) error
	GetQuestionByID(ctx context.Context, id uuid.UUID) (*entities.FeedbackQuestion, error)
	SubmitAnswer(ctx context.Context, answer *entities.FeedbackAnswer) error
	GetQuestionsWithAnswersByTeacher(ctx context.Context, teacherID uuid.UUID) ([]entities.FeedbackQuestion, error)
	GetFeedbackByTeacher(teacherID uuid.UUID) ([]entities.FeedbackQuestion, error)
}

type feedbackRepository struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &feedbackRepository{db}
}

func (r *feedbackRepository) CreateQuestion(ctx context.Context, question *entities.FeedbackQuestion) error {
	return r.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(question).Error
}



func (r *feedbackRepository) GetQuestionByID(ctx context.Context, id uuid.UUID) (*entities.FeedbackQuestion, error) {
	var question entities.FeedbackQuestion
	err := r.db.WithContext(ctx).Preload("Answers").First(&question, "id = ?", id).Error
	return &question, err
}

func (r *feedbackRepository) GetFeedbackByTeacher(teacherID uuid.UUID) ([]entities.FeedbackQuestion, error) {
	var feedbacks []entities.FeedbackQuestion

	// hanya ambil pertanyaan (tanpa preload answers)
	err := r.db.
		Where("created_by = ?", teacherID).
		Order("created_at DESC").
		Find(&feedbacks).Error

	if err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (r *feedbackRepository) SubmitAnswer(ctx context.Context, answer *entities.FeedbackAnswer) error {
	return r.db.WithContext(ctx).Create(answer).Error
}


func (r *feedbackRepository) GetQuestionsWithAnswersByTeacher(ctx context.Context, teacherID uuid.UUID) ([]entities.FeedbackQuestion, error) {
	var questions []entities.FeedbackQuestion
	err := r.db.WithContext(ctx).
		Preload("Answers").
		Preload("Answers.Student").
		Where("created_by = ?", teacherID).
		Find(&questions).Error
	return questions, err
}
