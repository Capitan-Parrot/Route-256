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

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE students;
-- +goose StatementEnd
