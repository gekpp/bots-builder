package questionnaire

import (
	"context"

	"github.com/google/uuid"
)

type dummy struct {
}

func NewDummy(qnrID uuid.UUID) Service {
	d := dummy{}
	return &d
}

// Start reset question asked to the user
// and returns welcome message and the first question.
func (d *dummy) Start(ctx context.Context, userID uuid.UUID) (StartResponse, error) {
	panic("not implemented") // TODO: Implement
}

// Answer validates, saves anser and returns next question to ask.
func (d *dummy) Answer(ctx context.Context, userID uuid.UUID, answer Answer) (AnswerResponse, error) {
	panic("not implemented") // TODO: Implement
}
