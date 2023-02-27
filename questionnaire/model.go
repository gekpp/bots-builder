package questionnaire

import (
	"errors"

	"github.com/google/uuid"
)

const invalidAnswerInfoMessage = "Пожалуйста, ответьте корректно на заданный вопрос 🧐"

var (
	errNoMoreQuestions = errors.New("no more questions")
	errInvalidAnswer   = errors.New("invalid answer")
)

type (
	questionKind      string
	userQuestionState string

	questionnaire struct {
		ID              uuid.UUID     `db:"id"`
		WelcomeMessage  string        `db:"welcome_message"`
		GoodbyeMessage  string        `db:"goodbye_message"`
		StartQuestionID uuid.NullUUID `db:"start_question_id"`
	}

	question struct {
		QuestionnaireID uuid.UUID      `db:"questionnaire_id"`
		ID              uuid.UUID      `db:"id"`
		Question        string         `db:"question"`
		Kind            questionKind   `db:"kind"`
		NextQuestionID  uuid.NullUUID  `db:"next_question_id"`
		AnswerOptions   []answerOption `db:"-"`
		RangeAnswer     rangeAnswer    `db:"-"`
	}

	answerOption struct {
		ID             uuid.UUID     `db:"id"`
		QuestionID     uuid.UUID     `db:"question_id"`
		Answer         string        `db:"answer"`
		NextQuestionID uuid.NullUUID `db:"next_question_id"`
		Rank           int           `db:"rank"`
	}

	userAnswers struct {
		ID              uuid.UUID         `db:"id"`
		UserID          uuid.UUID         `db:"user_id"`
		QuestionnaireID uuid.UUID         `db:"questionnaire_id"`
		QuestionID      uuid.UUID         `db:"question_id"`
		RawAnswer       string            `db:"raw_answer"`
		QuestionState   userQuestionState `db:"question_state"`
	}

	rangeAnswer struct {
		ID         uuid.UUID `db:"id"`
		QuestionID uuid.UUID `db:"question_id"`
		Minimum    int       `db:"minimum"`
		Maximum    int       `db:"maximum"`
	}
)

const (
	questionKindClose questionKind = "close"
	questionKindOpen  questionKind = "open"
	questionKindRange questionKind = "range"
)

const (
	answerStateAsked    userQuestionState = "asked"
	answerStateAnswered userQuestionState = "answered"
)
