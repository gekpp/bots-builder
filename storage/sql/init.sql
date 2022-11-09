CREATE TABLE users (
    id uuid NOT NULL primary key,
    user_name varchar(50)
);

CREATE TYPE question_type AS ENUM ('close', 'range', 'open');
CREATE TABLE questions (
    id uuid NOT NULL primary key,
    text varchar NOT NULL,
    kind question_type NOT NULL,
    next_question_id uuid REFERENCES questions(id),
    questionnaire_id uuid
);

CREATE TABLE questionnaires (
    id uuid NOT NULL PRIMARY KEY,
    name varchar(50),
    welcome_message varchar (200),
    start_question_id uuid NOT NULL REFERENCES questions(id)
);

ALTER TABLE questions
ADD FOREIGN KEY (questionnaire_id)
REFERENCES questionnaires(id);

CREATE TABLE question_answers (
    id uuid NOT NULL PRIMARY KEY,
    question_id uuid NOT NULL REFERENCES questions(id),
    text varchar NOT NULL,
    next_question_id uuid REFERENCES questions(id)
);

CREATE TABLE range_answers (
    id uuid NOT NULL PRIMARY KEY,
    question_id uuid REFERENCES question_answers(id),
    min int,
    max int
);

CREATE TYPE answer_state AS ENUM ('asked', 'answered');
CREATE TABLE answers (
    id uuid NOT NULL PRIMARY KEY,
    user_id uuid REFERENCES users(id),
    question_id uuid REFERENCES questions(id),
    question_answer_id uuid REFERENCES question_answers(id),
    range_answer_id uuid REFERENCES range_answers(id),
    raw_answer varchar,
    created_at date,
    state answer_state
);

