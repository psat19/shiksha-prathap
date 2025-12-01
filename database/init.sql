CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    hashed_password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255),
    age INT,
    phone VARCHAR(10)
);