
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users_task(
    id SERIAL PRIMARY KEY,
    user_id uuid REFERENCES users(id),
    is_completed BOOLEAN DEFAULT FALSE,
    archived_at TIMESTAMP WITH TIME ZONE
);