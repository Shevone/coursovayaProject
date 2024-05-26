CREATE TABLE lessons (
                         id SERIAL PRIMARY KEY,
                         name TEXT,
                         trainer_id INTEGER,
                         available_seats INTEGER,
                         description TEXT,
                         difficulty VARCHAR(8) CHECK (difficulty IN ('EASY', 'MEDIUM', 'HARD')),
                         startTime TIME,
                         day_of_week INTEGER CHECK (day_of_week BETWEEN 0 AND 6)
);

-- Создание таблицы student_lessons
CREATE TABLE student_lessons (
                                 lesson_id INTEGER REFERENCES lessons ON DELETE CASCADE,
                                 user_id INTEGER,
                                 PRIMARY KEY (lesson_id, user_id)
);
INSERT INTO lessons (name, trainer_id, available_seats, description, difficulty, startTime, day_of_week) VALUES
                                                                                                             ('Йога для начинающих', 2, 12, 'Расслабляющая йога для новичков', 'EASY', '09:00', 0), -- Воскресенье
                                                                                                             ('Силовая тренировка', 2, 10, 'Тренировка для всех групп мышц', 'MEDIUM', '10:30', 1), -- Понедельник
                                                                                                             ('Зумба', 2, 15, 'Танцевальная фитнес-программа', 'MEDIUM', '18:00', 2), -- Вторник
                                                                                                             ('Пилатес', 2, 8, 'Укрепление мышц и улучшение гибкости', 'EASY', '10:00', 3), -- Среда
                                                                                                             ('Фитнес для женщин', 2, 12, 'Комплексная тренировка для женщин', 'MEDIUM', '17:00', 4), -- Четверг
                                                                                                             ('Бокс для начинающих', 2, 10, 'Введение в основы бокса', 'MEDIUM', '19:00', 5), -- Пятница
                                                                                                             ('Растяжка и гибкость', 2, 15, 'Упражнения для улучшения гибкости', 'EASY', '11:00', 6); -- Суббота