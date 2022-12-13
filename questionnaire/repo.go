package questionnaire

import (
	"context"
	"database/sql"
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

func (r *repo) getQuestionnaire(ctx context.Context, qnrID uuid.UUID) (questionnaire, error) {
	res := questionnaire{}
	db := r.db.Unsafe()
	err := db.GetContext(ctx,
		&res,
		`SELECT * FROM questionnaires WHERE id=$1`,
		qnrID)
	if errors.Is(err, sql.ErrNoRows) {
		return questionnaire{}, ErrNotFound
	}
	if err != nil {
		return questionnaire{}, err
	}

	return res, nil
}

func (r *repo) getFirstQuestion(ctx context.Context, qnrID uuid.UUID) (question, error) {

	qnr, err := r.getQuestionnaire(ctx, qnrID)
	if err != nil {
		return question{}, fmt.Errorf("repo.getQuestionnaire: %w", err)
	}

	if !qnr.StartQuestionID.Valid {
		return question{}, ErrNoStartQuestion
	}

	res, err := r.getQuestion(ctx, qnr.StartQuestionID.UUID)
	if err != nil {
		return question{}, fmt.Errorf("repo.getQuestion: %w", err)
	}

	return res, nil
}

func (r *repo) getQuestionAnswerOptions(ctx context.Context, qID uuid.UUID) ([]answerOption, error) {
	db := r.db.Unsafe()
	res := []answerOption{}
	err := db.SelectContext(
		ctx,
		&res,
		"SELECT * FROM answer_options WHERE question_id=$1 ORDER BY rank ASC", qID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repo) getQuestionRangeAnswer(ctx context.Context, qID uuid.UUID) (rangeAnswer, error) {
	db := r.db.Unsafe()
	res := rangeAnswer{}
	err := db.GetContext(
		ctx,
		&res,
		"SELECT * FROM range_answer WHERE question_id=$1", qID)
	if errors.Is(err, sql.ErrNoRows) {
		return rangeAnswer{}, fmt.Errorf("range answer for question not found: %w", ErrNotFound)
	}
	if err != nil {
		return rangeAnswer{}, err
	}

	return res, nil
}

// GetLatestAskedQuestion looks up for latest asked question and returns it or
// errNotFound
func (r *repo) getLatestAskedQuestion(
	ctx context.Context,
	qnrID uuid.UUID,
	userID uuid.UUID) (question, error) {

	return question{}, errors.New("not implemented")
}

func (r *repo) getQuestion(
	ctx context.Context,
	id uuid.UUID) (question, error) {

	db := r.db.Unsafe()
	res := question{}
	err := db.GetContext(ctx, &res, "SELECT * FROM questions WHERE id=$1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return question{}, fmt.Errorf("question[id=%v] not found: %w", id, ErrNotFound)
	}

	switch res.Kind {
	case questionKindClose:
		opts, err := r.getQuestionAnswerOptions(ctx, res.ID)
		if err != nil {
			return question{}, fmt.Errorf("repo.getQuestionAnswerOptions: %w", err)
		}
		if len(opts) == 0 {
			return res, ErrNoAnswerOptions
		}
		res.AnswerOptions = opts
	case questionKindRange:
		rng, err := r.getQuestionRangeAnswer(ctx, res.ID)
		if err != nil {
			return question{}, fmt.Errorf("repo.getQuestionRangeAnswer: %w", err)
		}
		res.RangeAnswer = rng
	}
	return res, nil
}

// ыaveAskedQuestion saves question as asked
func (r *repo) saveAskedQuestion(
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
