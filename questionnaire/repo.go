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
		return rangeAnswer{}, fmt.Errorf("range answer for question id=%v not found: %w", qID, ErrNotFound)
	}
	if err != nil {
		return rangeAnswer{}, err
	}

	return res, nil
}

func (r *repo) getLatestQuestionWithState(
	ctx context.Context,
	qnrID uuid.UUID,
	userID uuid.UUID,
	state userQuestionState) (question, error) {
	var qID uuid.UUID

	err := r.db.GetContext(ctx, &qID,
		`SELECT question_id 
		FROM user_answers
		WHERE 
			questionnaire_id=$1 and user_id=$2
			and question_state=$3
		ORDER BY created_at DESC
		LIMIT 1`, qnrID, userID, state)

	if errors.Is(err, sql.ErrNoRows) {
		return question{}, fmt.Errorf("latest question with state=%v not found: %w",
			state, ErrNotFound)
	}
	if err != nil {
		return question{}, err
	}
	return r.getQuestion(ctx, qID)
}

// GetLatestAskedQuestion looks up for latest asked question and returns it or
// errNotFound
func (r *repo) getLatestAskedQuestion(
	ctx context.Context,
	qnrID uuid.UUID,
	userID uuid.UUID) (question, error) {

	return r.getLatestQuestionWithState(ctx, qnrID, userID, answerStateAsked)
}

// getLatestAnsweredQuestion looks up for latest answred question and returns it or
// errNotFound
func (r *repo) getLatestAnsweredQuestion(
	ctx context.Context,
	qnrID uuid.UUID,
	userID uuid.UUID) (question, error) {

	return r.getLatestQuestionWithState(ctx, qnrID, userID, answerStateAnswered)
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

// saveAskedQuestion saves question as asked
func (r *repo) saveAskedQuestion(
	ctx context.Context,
	qnrID uuid.UUID,
	userID uuid.UUID,
	questionID uuid.UUID) error {

	db := r.db.Unsafe()
	_, err := db.ExecContext(ctx, "INSERT INTO user_answers (user_id, questionnaire_id, question_id, question_state) VALUES ($1, $2, $3, $4)", userID, qnrID, questionID, answerStateAsked)
	if err != nil {
		return err
	}

	return nil
}

// saveAnswer saves answer and marks question as answered
func (r *repo) saveAnswer(
	ctx context.Context,
	userID uuid.UUID,
	questionID uuid.UUID,
	answer Answer) error {

	db := r.db.Unsafe()
	_, err := db.ExecContext(ctx,
		`UPDATE user_answers 
		SET 
			raw_answer=$1, question_state=$2
		WHERE 
			user_id=$3
			-- AND questionnaire_id=(SELECT questionnaire_id FROM questions WHERE id=$4)
			AND question_id=$4
			AND question_state = $5`,
		answer, answerStateAnswered, userID, questionID, answerStateAsked)
	if err != nil {
		return err
	}

	return nil
}
