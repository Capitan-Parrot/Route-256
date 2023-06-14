-- +goose Up
-- +goose StatementBegin
CREATE TABLE students
(
    id             BIGSERIAL
        CONSTRAINT students_pk PRIMARY KEY,
    name           varchar(256) not null,
    course_program varchar(256) not null,
    created_at     timestamp default now(),
    updated_at     timestamp default now()
);

CREATE TABLE tasks
(
    id          BIGSERIAL
        CONSTRAINT tasks_pk PRIMARY KEY,
    description text,
    deadline    timestamp,
    created_at  timestamp default now(),
    updated_at  timestamp default now()
);

CREATE TABLE solutions
(
    id         BIGSERIAL
        CONSTRAINT task_student_pk PRIMARY KEY,
    student_id bigint not null,
    task_id    bigint not null,
    status     varchar(255) default 'Waiting for approval',
    created_at timestamp    default now(),
    updated_at timestamp    default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE students;
DROP TABLE tasks;
DROP TABLE solutions;
-- +goose StatementEnd
