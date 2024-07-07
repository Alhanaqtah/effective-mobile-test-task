CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(50),
    surname VARCHAR(50),
    patronymic VARCHAR(50),
    address TEXT,
    passport_serie INTEGER UNIQUE,
    passport_number INTEGER UNIQUE
);

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),  
    title TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    done_at TIMESTAMP DEFAULT NULL
);