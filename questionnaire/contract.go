package questionnaire

import (
	"context"

	"github.com/google/uuid"
)

type (
	Message string
	Answer  string

	// Question represents a question to ask. It holds Text to ask and AnswerOptions
	// if exists.
	Question struct {
		Text          Message
		AnswerOptions []Answer
	}

	// StartResponse is a response on Start request. It holds Greeting text to
	// send to a user and a qustion to ask next.
	StartResponse struct {
		Welcome  Message
		Question Question
	}

	// AnswerResponse ...
	AnswerResponse struct {
		InfoMessage     Message
		QuestionMessage Question
	}

	// Service keeps tracking user's questions and answers. It is responsible to
	// store user's answers and to return next question to ask
	Service interface {
		// Start reset question asked to the user
		// and returns welcome message and the first question.
		Start(ctx context.Context, qnrID uuid.UUID, userID uuid.UUID) (StartResponse, error)

		// Answer validates, saves anser and returns next question to ask.
		Answer(ctx context.Context, qnrID uuid.UUID, userID uuid.UUID, answer Answer) (AnswerResponse, error)
	}
)
