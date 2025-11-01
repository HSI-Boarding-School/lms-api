package handlers

import (
	"api-shiners/api/handlers/dto"
	"api-shiners/pkg/feedback"
	"api-shiners/pkg/utils"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FeedbackController struct {
	service feedback.FeedbackService
}

func NewFeedbackController(service feedback.FeedbackService) *FeedbackController {
	return &FeedbackController{service: service}
}


// @Summary Create new feedback question
// @Description Admin membuat pertanyaan feedback baru
// @Tags Feedback
// @Accept json
// @Produce json
// @Param request body dto.CreateQuestionRequest true "Question payload"
// @Success 201 {object} entities.FeedbackQuestion
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/feedback/questions [post]
func (h *FeedbackController) CreateQuestion(c *fiber.Ctx) error {
	var req dto.CreateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request body", "BadRequest", nil)
	}

	if req.Question == "" {
		return utils.Error(c, http.StatusBadRequest, "Question is required", "ValidationError", nil)
	}

	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return utils.Error(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized", nil)
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, "Invalid user ID", "Unauthorized", nil)
	}

	ctx := context.Background()

	question, err := h.service.CreateQuestion(ctx, req.Question, userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "Failed to create question", "InternalServerError", nil)
	}

	return utils.Success(c, http.StatusCreated, "Feedback question created successfully", question, nil)
}


// @Summary Get all feedback questions
// @Description Mendapatkan semua pertanyaan feedback
// @Tags Feedback
// @Produce json
// @Success 200 {array} entities.FeedbackQuestion
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/feedback/questions [get]
func (h *FeedbackController) GetAllQuestions(c *fiber.Ctx) error {
	ctx := context.Background()
	questions, err := h.service.GetAllQuestions(ctx)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "Failed to get questions", "InternalServerError", nil)
	}

	return utils.Success(c, http.StatusOK, "Get all feedback questions successfully", questions, nil)
}


// @Summary Submit feedback answer
// @Description Mahasiswa mengirimkan jawaban feedback
// @Tags Feedback
// @Accept json
// @Produce json
// @Param request body dto.SubmitAnswerRequest true "Answer payload"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/feedback/answers [post]
func (h *FeedbackController) SubmitAnswer(c *fiber.Ctx) error {
	var req dto.SubmitAnswerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid request body", "BadRequest", nil)
	}

	if req.QuestionID == "" || req.Answer == "" {
		return utils.Error(c, http.StatusBadRequest, "QuestionID and Answer are required", "ValidationError", nil)
	}

	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return utils.Error(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized", nil)
	}

	studentID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, "Invalid user ID", "Unauthorized", nil)
	}

	questionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "Invalid question ID", "BadRequest", nil)
	}

	ctx := context.Background()
	if err := h.service.SubmitAnswer(ctx, questionID, studentID, req.Answer); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "Failed to submit answer", "InternalServerError", nil)
	}

	return utils.Success(c, http.StatusCreated, "Answer submitted successfully", nil, nil)
}



func (h *FeedbackController) GetQuestionsWithAnswersByTeacher(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return utils.Error(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized", nil)
	}

	teacherID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return utils.Error(c, http.StatusUnauthorized, "Invalid teacher ID", "Unauthorized", nil)
	}

	ctx := context.Background()
	questions, err := h.service.GetQuestionsWithAnswersByTeacher(ctx, teacherID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "Failed to get feedback data", "InternalServerError", nil)
	}

	return utils.Success(c, http.StatusOK, "Get feedback with answers successfully", questions, nil)
}
