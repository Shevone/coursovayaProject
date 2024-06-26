-- Создание типа данных для уровня сложности
CREATE TYPE difficulty_level AS ENUM ('EASY', 'MEDIUM', 'HARD');

-- Создание таблицы lessons
CREATE TABLE lessons (
                         id SERIAL PRIMARY KEY,
                         name TEXT,
                         trainer_id INTEGER,
                         available_seats INTEGER,
                         description TEXT,
                         difficulty difficulty_level,
                         startTime TIME,
                         day_of_week INTEGER CHECK (day_of_week BETWEEN 0 AND 6)
);

-- Создание таблицы student_lessons
CREATE TABLE student_lessons (
                                 lesson_id INTEGER REFERENCES lessons ON DELETE CASCADE,
                                 user_id INTEGER,
                                 PRIMARY KEY (lesson_id, user_id)
);

-- Вставка данных в lessons
INSERT INTO lessons (name, trainer_id, available_seats, description, difficulty, startTime, day_of_week) VALUES
                                                                                                             ('Йога для начинающих', 2, 12, 'Расслабляющая йога для новичков', 'EASY', '09:00', 0),
                                                                                                             ('Силовая тренировка', 2, 10, 'Тренировка для всех групп мышц', 'MEDIUM', '10:30', 1),
                                                                                                             ('Зумба', 2, 15, 'Танцевальная фитнес-программа', 'MEDIUM', '18:00', 2),
                                                                                                             ('Пилатес', 2, 8, 'Укрепление мышц и улучшение гибкости', 'EASY', '10:00', 3),
                                                                                                             ('Фитнес для женщин', 2, 12, 'Комплексная тренировка для женщин', 'MEDIUM', '17:00', 4),
                                                                                                             ('Бокс для начинающих', 2, 10, 'Введение в основы бокса', 'MEDIUM', '19:00', 5),
                                                                                                             ('Растяжка и гибкость', 2, 15, 'Упражнения для улучшения гибкости', 'EASY', '11:00', 6);