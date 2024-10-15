package migrations

const schema2 = `
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS conferences;
DROP TABLE IF EXISTS conference_members;
DROP TABLE IF EXISTS messages;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    last_login TIMESTAMP
);

CREATE TABLE conferences (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    last_message INTEGER
);

CREATE TABLE conference_members (
    conference_id INTEGER REFERENCES conferences(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (conference_id, user_id)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    conference_id INTEGER REFERENCES conferences(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    text TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_messages_conference_id ON messages(conference_id);
`
