package questionnaire

import "github.com/google/uuid"

type (
	questionKind string

	questionnaire struct {
		ID              uuid.UUID
		WelcomMessage   string
		GoodbyeMessage  string
		StartQuestionID uuid.UUID
	}

	question struct {
		QuestionnaireID uuid.UUID
		ID              uuid.UUID
		Text            string
		Kind            questionKind
		NextQuestionID  uuid.UUID
	}

	answerOption struct {
		ID             uuid.UUID
		QuestionID     uuid.UUID
		Text           string
		NextQuestionID uuid.UUID
		Rank           int
	}

	rangeAnswer struct {
		ID             uuid.UUID
		QuestionID     uuid.UUID
		Min            int
		Max            int
		NextQuestionID uuid.UUID
		Rank           int
	}
)

const (
	questionKindClose questionKind = "close"
	questionKindOpen  questionKind = "open"
	questionKindRange questionKind = "range"
)
