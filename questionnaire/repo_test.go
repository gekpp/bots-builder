package questionnaire

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/gekpp/bots-builder/internal/tests"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var shutdownFunc func()

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	// order matters
	binds := []string{
		fmt.Sprintf("%s/../scripts/migrations/20221210-init-schema.sql", workingDir),    // 1
		fmt.Sprintf("%s/../scripts/migrations/20221212-test-init-data.sql", workingDir), // 1
	}

	bindsMap := map[string]string{}

	names := []rune{'a', 'a', 'a'}
	k := len(names) - 1
	for _, s := range binds {
		bindsMap[s] = fmt.Sprintf("/docker-entrypoint-initdb.d/%s.sql", string(names))
		for i := k; i <= k; i-- {
			if names[i] < 'z' { // if panics increase names slice
				names[i]++
				break
			}
			if names[i] == 'z' {
				names[i] = 'a'
			}
		}
	}
	db, shutdownFunc = tests.ConnectTestContainers(bindsMap)
	defer shutdownFunc()

	os.Exit(m.Run())
}

func Test_repo_getQuestionnaire(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		qnrID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    questionnaire
		wantErr bool
	}{
		{
			name: "not_found",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:   context.Background(),
				qnrID: uuid.New(),
			},
			wantErr: true,
		},
		{
			name: "found",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx:   context.Background(),
				qnrID: uuid.MustParse("a2bdff81-3fd3-46d3-9348-c2b6946ecd9a"),
			},
			want: questionnaire{
				ID:              uuid.MustParse("a2bdff81-3fd3-46d3-9348-c2b6946ecd9a"),
				WelcomeMessage:  "Привет!",
				GoodbyeMessage:  "Пока!",
				StartQuestionID: uuid.NullUUID{},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repo{
				db: tt.fields.db,
			}
			got, err := r.getQuestionnaire(tt.args.ctx, tt.args.qnrID)
			if err != nil {
				t.Logf("Error: %+v", err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("repo.getQuestionnaire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repo.getQuestionnaire() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repo_getQuestion(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        question
		wantErr     bool
		expectedErr error
	}{
		{
			name: "not_found",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
			},
			wantErr:     true,
			expectedErr: ErrNotFound,
		},
		{
			name: "open_found",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("70ac2bd8-4501-40b7-b993-dd1d84f096b9"),
			},
			want: question{
				QuestionnaireID: uuid.MustParse("a2bdff81-3fd3-46d3-9348-c2b6946ecd9a"),
				ID:              uuid.MustParse("70ac2bd8-4501-40b7-b993-dd1d84f096b9"),
				Question:        "Как тебя зовут?",
				Kind:            questionKindOpen,
				NextQuestionID:  uuid.NullUUID{},
				AnswerOptions:   nil,
				RangeAnswer:     rangeAnswer{},
			}},
		{
			name: "close_found",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("fe3148b2-8743-4a03-ab13-9244d76d9152"),
			},
			want: question{
				QuestionnaireID: uuid.MustParse("a2bdff81-3fd3-46d3-9348-c2b6946ecd9a"),
				ID:              uuid.MustParse("fe3148b2-8743-4a03-ab13-9244d76d9152"),
				Question:        "В какое время суток Вы просыпаетесь?",
				Kind:            questionKindClose,
				NextQuestionID:  uuid.NullUUID{},
				AnswerOptions: []answerOption{
					{
						ID:             uuid.MustParse("20bc51bd-8b01-4e9b-a885-f8e3df1857bd"),
						QuestionID:     uuid.MustParse("fe3148b2-8743-4a03-ab13-9244d76d9152"),
						Answer:         "Утром",
						NextQuestionID: uuid.NullUUID{},
						Rank:           1,
					},
					{
						ID:             uuid.MustParse("3fa4f4e5-4e45-4017-b3f4-cd251b7b3227"),
						QuestionID:     uuid.MustParse("fe3148b2-8743-4a03-ab13-9244d76d9152"),
						Answer:         "Днём",
						NextQuestionID: uuid.NullUUID{},
						Rank:           2,
					},
					{
						ID:             uuid.MustParse("a756c0f9-0ece-4bff-b3a9-6820cee35335"),
						QuestionID:     uuid.MustParse("fe3148b2-8743-4a03-ab13-9244d76d9152"),
						Answer:         "Вечером",
						NextQuestionID: uuid.NullUUID{},
						Rank:           3,
					},
					{
						ID:             uuid.MustParse("d35f2fdd-9bd8-4cb6-be32-f6e56a0ec450"),
						QuestionID:     uuid.MustParse("fe3148b2-8743-4a03-ab13-9244d76d9152"),
						Answer:         "Ночью",
						NextQuestionID: uuid.NullUUID{},
						Rank:           4,
					},
				},
				RangeAnswer: rangeAnswer{},
			},
		},
		{
			name: "close_no_answer",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("8b6c8fb7-a709-4a02-b51e-da23dfd4a413"),
			},
			want:        question{},
			wantErr:     true,
			expectedErr: ErrNoAnswerOptions,
		},
		{
			name: "range_found",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("ebecfd9d-a114-4ef8-a105-3951c0db20d0"),
			},
			want: question{
				QuestionnaireID: uuid.MustParse("a2bdff81-3fd3-46d3-9348-c2b6946ecd9a"),
				ID:              uuid.MustParse("ebecfd9d-a114-4ef8-a105-3951c0db20d0"),
				Question:        "Сколько раз в день Вы чистите зубы?",
				Kind:            questionKindRange,
				NextQuestionID:  uuid.NullUUID{},
				RangeAnswer: rangeAnswer{
					ID:         uuid.MustParse("e4eddd95-c504-483a-bed3-5e3eb907e113"),
					QuestionID: uuid.MustParse("ebecfd9d-a114-4ef8-a105-3951c0db20d0"),
					Minimum:    0,
					Maximum:    20,
				},
			},
		},
		{
			name: "range_no_answer",
			fields: fields{
				db: sqlx.NewDb(db, "postgres"),
			},
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("64379850-defa-4eb6-96b2-36cb3f22bd0a"),
			},
			want:        question{},
			wantErr:     true,
			expectedErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repo{
				db: tt.fields.db,
			}
			got, err := r.getQuestion(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
