-- Создание базы данных
CREATE DATABASE database_for_news;

-- Подключение к базе данных
\c database_for_news;

-- Создание таблицы posts
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    pub_time TIMESTAMP NOT NULL,
    link TEXT NOT NULL UNIQUE
);