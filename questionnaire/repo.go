package questionnaire

import (
	"context"
	"errors"
	"fmt"

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

func (r *repo) getQuestionnaire(ctx context.Context, qnrID uuid.UUID) (questionnaire, error) {
	return questionnaire{}, errors.New("not implemented")
}

func (r *repo) getFirstQuestion(ctx context.Context, qnrID uuid.UUID) (question, error) {
	res := question{}

	if res.Kind == questionKindClose {
		opts, err := r.getQuestionAnswerOptions(ctx, res.ID)
		if err != nil {
			return res, fmt.Errorf("repo.getQuestionAnswerOptions: %v", err)
		}
		if len(opts) == 0 {
			return res, errNoAnswerOptions
		}
	}

	return res, errors.New("not implemented")
}

func (r *repo) getQuestionAnswerOptions(ctx context.Context, qID uuid.UUID) ([]answerOption, error) {
	return []answerOption{}, errors.New("not implemented")
}

// getNextQuestionByRank returns the question Q of the questionnaire[id=qnrID]
// with rank greather than
func (r *repo) getNextQuestionByRank(ctx context.Context, qnrID uuid.UUID, rank int) (question, error) {
	return question{}, errors.New("not implemented")
}

// GetLatestAskedQuestion looks up for latest asked question and returns it or
// errNotFound
func (r *repo) getLatestAskedQuestion(
	ctx context.Context,
	qnrID uuid.UUID,
	userID uuid.UUID) (question, error) {

	return question{}, errors.New("not implemented")
}

func (r *repo) GetQuestion(
	ctx context.Context,
	id uuid.UUID) (question, error) {
	return question{}, errors.New("not implemented")
}

// SaveAskedQuestion saves question as asked
func (r *repo) SaveAskedQuestion(
	ctx context.Context,
	qnrID uuid.UUID,
	userID uuid.UUID,
	questionID uuid.UUID) error {

	return errors.New("not implemented")
}

// saveAnswer saves answer and marks question as answered
func (r *repo) saveAnswer(
	ctx context.Context,
	qnrID uuid.UUID,
	userID uuid.UUID,
	answer Answer) error {

	return errors.New("not implemented")
}
