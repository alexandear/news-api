CREATE TABLE IF NOT EXISTS posts (
    id uuid PRIMARY KEY,
    title VARCHAR (50) NOT NULL,
    content TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
