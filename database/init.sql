CREATE DATABASE genshinquiz;

\c genshinquiz;

CREATE TABLE IF NOT EXISTS quizzes (
  id SERIAL PRIMARY KEY,
  question VARCHAR(255) NOT NULL,
  answer VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  email VARCHAR(100)
);

INSERT INTO users (name, email) VALUES ('John Doe', 'john.doe@example.com');