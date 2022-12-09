package questionnaire

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type service struct {
}

// New creates and returns new service instance
func New(db *sqlx.DB, qnrID uuid.UUID) *service {
	s := service{}
	return &s
}

// Start reset question asked and returns welcome message and the first question
func (s *service) Start(
	ctx context.Context,
	userID uuid.UUID) (StartResponse, error) {

	return StartResponse{}, errors.New("not implemented")
}

// Answer validates, saves and returns next question to ask
func (s *service) Answer(
	ctx context.Context,
	userID uuid.UUID,
	answer Answer) (AnswerResponse, error) {

	return AnswerResponse{}, errors.New("not implemented")
}
