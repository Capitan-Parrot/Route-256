---генерация данных
INSERT INTO students(name, course_program)
SELECT md5(random()::varchar(256)), md5(random()::varchar(256))
FROM generate_series(1, 100);

INSERT INTO tasks(description, deadline)
SELECT (array['Test', 'Project', 'Homework'])[floor(random() * 3 + 1)], NOW() + (random() * (interval '1 week'))
FROM generate_series(1, 100);

INSERT INTO solutions(student_id, task_id, status, created_at, updated_at)
SELECT floor(random() * 100 + 1), floor(random() * 100 + 1),
       (array['Waiting for approval', 'OK', 'Not OK'])[floor(random() * 3 + 1)],
       NOW() - '1 week'::interval - (random() * (interval '1 week')), NOW() + (random() * (interval '1 week'))
FROM generate_series(1, 100);