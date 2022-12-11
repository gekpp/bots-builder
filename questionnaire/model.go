package questionnaire

import (
	"errors"

	"github.com/google/uuid"
)

const invalidAnswerInfoMessage = "Пожалуйста, ответьте одним из предложенных вариантов ответа или значением из заданного интервала."

var (
	errNoAnswerOptions = errors.New("no answer options found")
	errNoMoreQuestions = errors.New("no more questions")
	errInvalidAnswer   = errors.New("invalid answer")
)

type (
	questionKind string

	questionnaire struct {
		ID              uuid.UUID
		WelcomeMessage  string
		GoodbyeMessage  string
		StartQuestionID uuid.UUID
	}

	question struct {
		QuestionnaireID uuid.UUID
		ID              uuid.UUID
		Question        string
		Kind            questionKind
		NextQuestionID  uuid.NullUUID
		AnswerOptions   []answerOption
		RangeAnswer     rangeAnswer
		Rank            int
	}

	answerOption struct {
		ID             uuid.UUID
		QuestionID     uuid.UUID
		Answer         string
		NextQuestionID uuid.NullUUID
		Rank           int
	}

	rangeAnswer struct {
		ID             uuid.UUID
		QuestionID     uuid.UUID
		Minimum        int
		Maximum        int
		NextQuestionID uuid.UUID
		Rank           int
	}
)

const (
	questionKindClose questionKind = "close"
	questionKindOpen  questionKind = "open"
	questionKindRange questionKind = "range"
)
