INSERT INTO
    public.questionnaires (
        id,
        name,
        start_question_id,
        welcome_message,
        goodbye_message,
        created_at,
        updated_at
    )
VALUES
    (
        'a2bdff81-3fd3-46d3-9348-c2b6946ecd9a',
        'Epileplsy',
        NULL,
        'Привет!',
        'Пока!',
        '2022-12-13 00:08:10.000000',
        '2022-12-12 22:08:15.445665'
    );

INSERT INTO
    public.questions (
        id,
        question,
        kind,
        next_question_id,
        questionnaire_id,
        created_at,
        updated_at
    )
VALUES
    (
        '70ac2bd8-4501-40b7-b993-dd1d84f096b9',
        'Как тебя зовут?',
        'open',
        null,
        'a2bdff81-3fd3-46d3-9348-c2b6946ecd9a',
        '2022-12-13 00:21:33.000000',
        '2022-12-12 22:22:09.565755'
    ),
    (
        'fe3148b2-8743-4a03-ab13-9244d76d9152',
        'В какое время суток Вы просыпаетесь?',
        'close',
        null,
        'a2bdff81-3fd3-46d3-9348-c2b6946ecd9a',
        '2022-12-12 22:37:17.810614',
        '2022-12-12 22:37:17.810614'
    ),
    (
        'ebecfd9d-a114-4ef8-a105-3951c0db20d0',
        'Сколько раз в день Вы чистите зубы?',
        'range',
        null,
        'a2bdff81-3fd3-46d3-9348-c2b6946ecd9a',
        '2022-12-12 22:37:17.810614',
        '2022-12-12 22:37:17.810614'
    ),
    (
        '64379850-defa-4eb6-96b2-36cb3f22bd0a',
        'Ранжированный вопрос без вариантов. Ошибка!',
        'range',
        null,
        'a2bdff81-3fd3-46d3-9348-c2b6946ecd9a',
        '2022-12-12 22:37:17.810614',
        '2022-12-12 22:37:17.810614'
    ),
    (
        '8b6c8fb7-a709-4a02-b51e-da23dfd4a413',
        'Закрытый вопрос без вариантов - Ошибка!',
        'close',
        null,
        'a2bdff81-3fd3-46d3-9348-c2b6946ecd9a',
        '2022-12-12 22:37:17.810614',
        '2022-12-12 22:37:17.810614'
    );

INSERT INTO
    public.answer_options (
        id,
        question_id,
        answer,
        rank,
        next_question_id,
        created_at,
        updated_at
    )
VALUES
    (
        '20bc51bd-8b01-4e9b-a885-f8e3df1857bd',
        'fe3148b2-8743-4a03-ab13-9244d76d9152',
        'Утром',
        1,
        null,
        '2022-12-12 22:38:36.593738',
        '2022-12-12 22:38:36.593738'
    ),
    (
        '3fa4f4e5-4e45-4017-b3f4-cd251b7b3227',
        'fe3148b2-8743-4a03-ab13-9244d76d9152',
        'Днём',
        2,
        null,
        '2022-12-12 22:38:36.593738',
        '2022-12-12 22:38:36.593738'
    ),
    (
        'a756c0f9-0ece-4bff-b3a9-6820cee35335',
        'fe3148b2-8743-4a03-ab13-9244d76d9152',
        'Вечером',
        3,
        null,
        '2022-12-12 22:38:36.593738',
        '2022-12-12 22:38:36.593738'
    ),
    (
        'd35f2fdd-9bd8-4cb6-be32-f6e56a0ec450',
        'fe3148b2-8743-4a03-ab13-9244d76d9152',
        'Ночью',
        4,
        null,
        '2022-12-12 22:38:36.593738',
        '2022-12-12 22:38:36.593738'
    );

INSERT INTO
    public.range_answer (
        id,
        question_id,
        minimum,
        maximum,
        created_at,
        updated_at
    )
VALUES
    (
        'e4eddd95-c504-483a-bed3-5e3eb907e113',
        'ebecfd9d-a114-4ef8-a105-3951c0db20d0',
        0,
        20,
        '2022-12-12 22:46:28.126594',
        '2022-12-12 22:46:28.126594'
    );

INSERT INTO public.user_answers (
    id, 
    user_id, 
    questionnaire_id, 
    question_id, 
    raw_answer, 
    created_at, 
    question_state, 
    updated_at) 
VALUES (
    '4c617a80-7c43-46ed-bc74-811fb07ed6f2', 
    'fe3148b2-8743-4a03-ab13-9244d76d9152', 
    'a2bdff81-3fd3-46d3-9348-c2b6946ecd9a', 
    '70ac2bd8-4501-40b7-b993-dd1d84f096b9', 
    '', 
    '2022-12-15 16:49:26.726834', 
    'asked', 
    '2022-12-15 16:49:26.726834'
    );
