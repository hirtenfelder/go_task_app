CREATE TABLE tasks
(
    task_pk     SERIAL PRIMARY KEY,
    name        TEXT    NOT NULL,
    description TEXT NULL,
    due_date    TEXT NULL,
    done        BOOLEAN NOT NULL DEFAULT FALSE
);