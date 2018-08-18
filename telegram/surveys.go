package telegram

type SurveyConfig struct {
	StartMessage    string
	Questions       []Question
	CompleteMessage string
}

type Question struct {
	Text string
	//If true user's phone number will be requested
	RequestContact bool
	//If true user's phone number will be requested unless user set up his username. If both are true contact will be requested
	//RequestContactUnlessUsername bool
}

func (q Question) String() string {
	return q.Text
}

type surveyOption func(s *SurveyConfig)

func WithQuestion(q Question) surveyOption {
	return func(s *SurveyConfig) {
		s.Questions = append(s.Questions, q)
	}
}

func WithStartMessage(msg string) surveyOption {
	return func(s *SurveyConfig) {
		s.StartMessage = msg
	}
}

func WithCompleteMessage(msg string) surveyOption {
	return func(s *SurveyConfig) {
		s.CompleteMessage = msg
	}
}

func NewSurveyConfig(question Question, opts ...surveyOption) *SurveyConfig {
	s := SurveyConfig{Questions: []Question{question}}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}
