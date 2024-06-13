-- Создание типа данных для ролей
CREATE TYPE user_role AS ENUM ('New', 'User', 'Trainer', 'Admin');

-- Создание таблицы users
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email TEXT UNIQUE NOT NULL,
                       pass_hash TEXT NOT NULL,
                       name TEXT NOT NULL,
                       surname TEXT NOT NULL,
                       patronymic TEXT,
                       phone_number TEXT NOT NULL,
                       role user_role NOT NULL
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
                                                                                            (1, 'admin', '$2a$10$EzojjVsI8oaj4wQTG7aCbuqRkFH02HFNii/m8vtFfOfZd2JvamsdO', 'Администратор', 'Админов', 'Админович', '938129084', 'Admin'),
                                                                                            (2, 'trainer', '$2a$10$VTZeyyU6vdNNg0DCbxzw9.p08lN6jartw6Qx4f2jYUAtbnRurO3N.', 'Тренер', 'Тренеров', 'Тренерович', '893894', 'Trainer'),
                                                                                            (3, 'user', '$2a$10$VTZeyyU6vdNNg0DCbxzw9.p08lN6jartw6Qx4f2jYUAtbnRurO3N.', 'Клиент', 'Клиентов', '', '893894', 'User'),
                                                                                            (4, 'newUser', '$2a$10$VTZeyyU6vdNNg0DCbxzw9.p08lN6jartw6Qx4f2jYUAtbnRurO3N.', 'Новый', 'Новый', 'Новый', '893894', 'New');

-- Вставка данных в admin_level
INSERT INTO admin_level (admin_id, admin_level) VALUES (1, 1);