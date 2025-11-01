package routes

import (
	"api-shiners/api/handlers"
	"api-shiners/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func FeedbackRoutes(app *fiber.App, feedbackController *handlers.FeedbackController) {
	api := app.Group("/api")

	api.Post("/feedback/questions",middleware.TeacherMiddleware, feedbackController.CreateQuestion)
	api.Get("/feedback/questions", middleware.AuthMiddleware, feedbackController.GetAllQuestions)
	
	api.Post("/feedback/answers", middleware.AuthMiddleware, feedbackController.SubmitAnswer)

	api.Get("/feedback/teacher", middleware.TeacherMiddleware, feedbackController.GetQuestionsWithAnswersByTeacher)
}


