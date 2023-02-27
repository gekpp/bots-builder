package questionnaire

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type service struct {
	r     *repo
	qnrID uuid.UUID
}

// New creates and returns new service instance
func New(db *sqlx.DB, qnrID uuid.UUID) *service {
	s := service{
		r:     &repo{db: db},
		qnrID: qnrID,
	}
	return &s
}

// Start reset question asked and returns welcome message and the first question
func (s *service) Start(
	ctx context.Context,
	userID uuid.UUID) (StartResponse, error) {

	qnr, err := s.r.getQuestionnaire(ctx, s.qnrID)
	if err != nil {
		return StartResponse{}, fmt.Errorf("could not get questionnaire: repo.getQuestionnaire: %w", err)
	}

	q, err := s.r.getFirstQuestion(ctx, s.qnrID)
	if err != nil {
		return StartResponse{}, fmt.Errorf("could not get first question: repo.getFirstQuestion: %w", err)
	}

	if err := s.r.saveAskedQuestion(ctx, s.qnrID, userID, q.ID); err != nil {
		return StartResponse{}, fmt.Errorf("could not save asked question: repo.SaveAskedQuestion: %w", err)
	}

	return StartResponse{
		Welcome: Message(qnr.WelcomeMessage),
		Question: Question{
			Text:          Message(q.Question),
			AnswerOptions: getAnswerOptionsIfRequired(q),
		},
	}, nil
}

// Answer validates, saves and returns next question to ask
func (s *service) Answer(
	ctx context.Context,
	userID uuid.UUID,
	answer Answer) (AnswerResponse, error) {

	qnr, err := s.r.getQuestionnaire(ctx, s.qnrID)
	if err != nil {
		return AnswerResponse{}, fmt.Errorf("could not get questionnaire id=%v: repo.getQuestionnaire: %v",
			s.qnrID, err)
	}

	latestQuestion, err := s.r.getLatestAskedQuestion(ctx, s.qnrID, userID)
	if errors.Is(err, ErrNotFound) {
		// FIXME: use the answer of the latest answered question to determine the next question
		// ignore answer from arguments.
		latestAnsweredQuestion, err := s.r.getLatestAnsweredQuestion(ctx, s.qnrID, userID)
		if err != nil {
			return AnswerResponse{}, fmt.Errorf("could not get latest answered question, qnrID=%v, userID=%v: repo.getLatestAskedQuestion: %v",
				s.qnrID, userID, err)
		}
		return s.askNextQuestion(ctx, qnr, userID, latestAnsweredQuestion, answer)
	}
	if err != nil {
		return AnswerResponse{}, fmt.Errorf("could not get latest asked question, qnrID=%v, userID=%v: repo.getLatestAskedQuestion: %v",
			s.qnrID, userID, err)
	}

	if err := validateAnswer(ctx, latestQuestion, answer); err != nil {
		return AnswerResponse{
			Info: invalidAnswerInfoMessage,
			Question: Question{
				Text:          Message(latestQuestion.Question),
				AnswerOptions: getAnswerOptionsIfRequired(latestQuestion),
			},
		}, fmt.Errorf("%w: userID=%v, questionID=%v, answer=%q", err, userID, latestQuestion.ID, answer)
	}

	if err := s.r.saveAnswer(ctx, userID, latestQuestion.ID, answer); err != nil {
		return AnswerResponse{}, fmt.Errorf("could not save answer: repo.saveAnswer: %v", err)
	}

	return s.askNextQuestion(ctx, qnr, userID, latestQuestion, answer)
}

func (s *service) askNextQuestion(ctx context.Context,
	qnr questionnaire,
	userID uuid.UUID,
	latestQuestion question,
	answer Answer) (AnswerResponse, error) {

	nextQ, err := s.getNextQuestion(ctx, latestQuestion, answer)
	switch {
	case err == nil:
		break
	case errors.Is(err, errNoMoreQuestions):
		return AnswerResponse{
			Info: Message(qnr.GoodbyeMessage),
		}, nil
	case errors.Is(err, errInvalidAnswer):
		return s.askQuestion(ctx, qnr, userID, latestQuestion)
	default:
		return AnswerResponse{}, fmt.Errorf("could not get next question: getNextQuestion: %v", err)
	}

	return s.askQuestion(ctx, qnr, userID, nextQ)

}

func (s *service) askQuestion(ctx context.Context,
	qnr questionnaire,
	userID uuid.UUID,
	q question) (AnswerResponse, error) {

	if err := s.r.saveAskedQuestion(ctx, qnr.ID, userID, q.ID); err != nil {
		return AnswerResponse{}, fmt.Errorf("could not save asked question: repo.saveAskedQuestion: %v", err)
	}

	return AnswerResponse{
		Question: Question{
			Text:          Message(q.Question),
			AnswerOptions: getAnswerOptionsIfRequired(q),
		},
	}, nil
}

func (s *service) getNextQuestion(ctx context.Context, q question, a Answer) (question, error) {

	nextQID := uuid.NullUUID{}

	switch q.Kind {
	case questionKindOpen:
		nextQID = q.NextQuestionID
	case questionKindRange:
		nextQID = q.NextQuestionID
	case questionKindClose:
		found := false
		for _, op := range q.AnswerOptions {
			if a == Answer(op.Answer) {
				found = true
				nextQID = op.NextQuestionID
				if !nextQID.Valid {
					nextQID = q.NextQuestionID
				}
			}
		}
		if !found {
			return question{}, errInvalidAnswer
		}
	default:
		return question{}, errors.New("invalid question kind")
	}

	if !nextQID.Valid {
		return question{}, errNoMoreQuestions
	}

	q, err := s.r.getQuestion(ctx, nextQID.UUID)
	if err != nil {
		return question{}, fmt.Errorf("repo.GetQuestion: %v", err)
	}
	return q, nil
}

func getAnswerOptionsIfRequired(q question) []Answer {

	if q.Kind != questionKindClose {
		return nil
	}

	res := make([]Answer, 0, len(q.AnswerOptions))
	for _, o := range q.AnswerOptions {
		res = append(res, Answer(o.Answer))
	}
	return res
}

func validateAnswer(ctx context.Context, q question, a Answer) error {
	switch q.Kind {
	case questionKindOpen:
		return nil
	case questionKindClose:
		for _, o := range q.AnswerOptions {
			if a == Answer(o.Answer) {
				return nil
			}
		}
		return errInvalidAnswer
	case questionKindRange:
		return validateRangeAnswer(a, q.RangeAnswer.Minimum, q.RangeAnswer.Maximum)
	default:
		return fmt.Errorf("unknown question kind %v", q.Kind)
	}
}

func validateRangeAnswer(a Answer, min, max int) error {
	i, err := strconv.Atoi(string(a))
	if err != nil {
		return fmt.Errorf("%w: could not convert answer to int: %v", errInvalidAnswer, err)
	}

	if min <= i && i <= max {
		return nil
	}

	return fmt.Errorf("%w: answer value %v is our of range [%v, %v]", errInvalidAnswer, i, min, max)
}
