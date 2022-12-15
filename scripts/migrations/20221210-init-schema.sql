create table users (
    id uuid default gen_random_uuid() not null constraint users_pk primary key,
    telegram_user_id text,
    telegram_username text,
    created_at timestamp without time zone default now() not null,
    updated_at timestamp without time zone default now() not null
);

create unique index users_telegram_user_id_uindex on users (telegram_user_id);

create table questionnaires (
    id uuid default gen_random_uuid() not null constraint questionnaires_pk primary key,
    name text not null,
    welcome_message text not null,
    goodbye_message text not null,
    start_question_id uuid,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create type question_kind AS ENUM ('open', 'close', 'range');

create table questions (
    id uuid default gen_random_uuid() not null constraint questions_pk primary key,
    question text not null,
    kind question_kind not null,
    next_question_id uuid,
    questionnaire_id uuid not null constraint questions_questionnaires_id_fk references questionnaires,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create table range_answer (
    id uuid default gen_random_uuid() not null constraint range_answer_pk primary key,
    question_id uuid not null constraint range_answer_questions_id_fk references questions,
    minimum int not null,
    maximum int not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create unique index range_answer_question_id_uindex on range_answer (question_id);

create table answer_options (
    id uuid default gen_random_uuid() not null constraint answer_options_pk primary key,
    question_id uuid not null constraint answer_options_questions_id_fk references questions,
    answer text not null,
    rank integer not null,
    next_question_id uuid,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create unique index answer_options_question_id_rank_uindex on answer_options (question_id, rank);

create type user_question_state as enum ('asked', 'answered');

create table user_answers (
    id uuid default gen_random_uuid() not null constraint user_answers_pk primary key,
    user_id uuid not null,
    questionnaire_id uuid not null,
    question_id uuid not null,
    raw_answer text default '',
    created_at timestamp without time zone default now() not null,
    question_state user_question_state not null,
    updated_at timestamp without time zone default now() not null
);

create index user_answers_user_id_questionnaire_id_created_at_index on user_answers (
    user_id asc,
    questionnaire_id asc,
    created_at desc
)
WHERE
    question_state = 'asked';