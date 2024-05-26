-- Создание таблицы lessons
CREATE TABLE lessons (
                         id SERIAL PRIMARY KEY,
                         title TEXT,
                         trainer_id INTEGER,
                         available_seats INTEGER,
                         description TEXT,
                         difficult VARCHAR(8),
                         date_and_time TIMESTAMP WITHOUT TIME ZONE,
                         is_complete BOOLEAN DEFAULT FALSE,
                         CONSTRAINT difficult_check CHECK (difficult IN ('EASY', 'HARD', 'MEDIUM'))
);

-- Создание триггера для обновления is_complete
CREATE OR REPLACE FUNCTION update_is_complete()
RETURNS TRIGGER AS $$
BEGIN
UPDATE lessons
SET is_complete = CASE
                      WHEN NEW.available_seats = 0 THEN TRUE
                      ELSE FALSE
    END
WHERE id = NEW.id;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_is_complete_trigger
    AFTER UPDATE OF available_seats ON lessons
    FOR EACH ROW EXECUTE PROCEDURE update_is_complete();

-- Создание таблицы student_lessons
CREATE TABLE student_lessons (
                                 user_id INTEGER,
                                 lesson_id INTEGER REFERENCES lessons ON DELETE CASCADE,
                                 PRIMARY KEY (user_id, lesson_id)
);

-- Вставка данных в lessons
INSERT INTO lessons (id, title, trainer_id, available_seats, description, difficult, date_and_time, is_complete) VALUES
                                                                                                                     (1, 'incididunt irure aliqua ex', 69213661, 2, 'ex adipisicing', 'MEDIUM', '2023-12-01 12:00:00', FALSE),
                                                                                                                     (2, 'dkkdkdkdk', 3, 1, 'ьdlfdksl', 'HARD', '2023-12-02 14:00:00', FALSE),
                                                                                                                     (3, 'anim aliquip laborum elit in', 661174, 10, 'in Ut ad nostrud qui', 'EASY', '2023-12-03 16:00:00', FALSE),
                                                                                                                     (4, 'laborum consequat', 51869590, 8887, 'ex', 'HARD', '2023-12-04 18:00:00', FALSE),
                                                                                                                     (5, 'Duis velit est', 661174, 0, 'aliqua eu Duis', 'EASY', '2023-12-05 20:00:00', TRUE);