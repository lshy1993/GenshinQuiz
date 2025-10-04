-- This file ensures the database is properly initialized
-- It's automatically run by PostgreSQL when the container starts for the first time

-- Create the main database if it doesn't exist
SELECT 'CREATE DATABASE genshin_quiz'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'genshin_quiz')\gexec

-- Connect to the genshin_quiz database
\c genshin_quiz;

-- Create any initial extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
SET timezone = 'UTC';