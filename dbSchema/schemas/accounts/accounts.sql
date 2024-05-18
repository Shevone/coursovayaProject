-- Создание таблицы users
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email TEXT UNIQUE NOT NULL,
                       pass_hash TEXT NOT NULL,
                       name TEXT NOT NULL,
                       surname TEXT NOT NULL,
                       patronymic TEXT,
                       phone_number TEXT NOT NULL,
                       role TEXT NOT NULL
);

CREATE INDEX idx_email ON users (email);

-- Создание таблицы admin_level
CREATE TABLE admin_level (
                             admin_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                             admin_level INTEGER NOT NULL,
                             PRIMARY KEY (admin_id)
);

-- Вставка данных в users
INSERT INTO users (id, email, pass_hash, name, surname, patronymic, phone_number, role) VALUES
                                                                                            (1, 'kola2003@yandex.ru', '$2a$10$EzojjVsI8oaj4wQTG7aCbuqRkFH02HFNii/m8vtFfOfZd2JvamsdO', 'Nickolay', 'Ryabov', 'Dmitr', '938129084', 'Admin'),
                                                                                            (2, 'trainer@yandex.ru', '$2a$10$VTZeyyU6vdNNg0DCbxzw9.p08lN6jartw6Qx4f2jYUAtbnRurO3N.', 'krutoychel', 'Ryabov', '', '893894', 'Trainer'),
                                                                                            (3, 'user@yandex.ru', '$2a$10$VTZeyyU6vdNNg0DCbxzw9.p08lN6jartw6Qx4f2jYUAtbnRurO3N.', 'user', 'user', '', '893894', 'User');

-- Вставка данных в admin_level
INSERT INTO admin_level (admin_id, admin_level) VALUES
    (1, 10);