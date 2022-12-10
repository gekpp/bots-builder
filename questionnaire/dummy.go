package questionnaire

import (
	"context"

	"github.com/google/uuid"
)

type dummy struct {
}

var (
	firstQuestion = Question{
		Text:          "Вводное сообщение. Информированное согласие. Текст на согласовывании с ЛПУ.",
		AnswerOptions: []Answer{"Далее"},
	}

	secondQuestion = Question{
		Text: "За последние 4 недели общее состояние здоровья Вашего ребенка",
		AnswerOptions: []Answer{
			"Выраженное улучшение",
			"Значительное улучшение",
			"Незначительное улучшение",
			"Изменений нет",
			"Незначительное ухудшение",
			"Значительное ухудшение",
			"Выраженное ухудшение",
		},
	}
)

func NewDummy(qnrID uuid.UUID) Service {
	d := dummy{}
	return &d
}

// Start reset question asked to the user
// and returns welcome message and the first question.
func (d *dummy) Start(ctx context.Context, userID uuid.UUID) (StartResponse, error) {
	return StartResponse{
		Welcome: "Привет! Пожалуйста, пройдите опрос. " +
			"Опрос состоит из нескольких вопросов. " +
			"Пожалуйста, отвечайте на вопросы по одному. " +
			"По окончании опроса, будет финальное сообщение.",
		Question: firstQuestion,
	}, nil
}

// Answer validates if answer equal to "Далее" and return next question or ask
// the first one again.
// This method emulates service behaviour for the first question only.
func (d *dummy) Answer(ctx context.Context, userID uuid.UUID, answer Answer) (AnswerResponse, error) {
	if !validAnswer(answer) {
		return AnswerResponse{
			Info:     "Пожалуйста, выберите один из предложенных ответов",
			Question: firstQuestion,
		}, nil
	}
	return AnswerResponse{
		Question: secondQuestion,
	}, nil
}

func validAnswer(answer Answer) bool {
	return answer == "Далее"
}
