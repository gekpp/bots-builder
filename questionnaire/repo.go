package questionnaire

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	repo struct {
		db *sqlx.DB
	}
)

var (
	errNotFound error = errors.New("not found")
)

func (r *repo) GetWelcomeMessage(ctx context.Context, qnrID uuid.UUID) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (r *repo) GetInitialQuestion(ctx context.Context, qnrID uuid.UUID) (interface{}, error) {
	return nil, errors.New("not implemented")
}

// GetLatestAskedQuestion looks up for latest asked question and returns it or
// errNotFound
func (r *repo) GetLatestAskedQuestion(
	ctx context.Context,
	qnrID uuid.UUID,
	userID uuid.UUID) (interface{}, error) {

	return nil, errors.New("not implemented")
}

// SaveAskedQuestion saves question as asked
func (r *repo) SaveAskedQuestion(
	ctx context.Context,
	qnrID uuid.UUID,
	userID uuid.UUID,
	questionID uuid.UUID) (interface{}, error) {

	return nil, errors.New("not implemented")
}

// SaveAnswer saves answer and marks question as answered
func (r *repo) SaveAnswer(
	ctx context.Context,
	qnrID uuid.UUID,
	userID uuid.UUID,
	answer Answer) (interface{}, error) {

	return nil, errors.New("not implemented")
}
