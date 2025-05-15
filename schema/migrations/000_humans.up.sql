CREATE TABLE IF NOT EXISTS humans (
    ID SERIAL,
    name VARCHAR(50),
    surname VARCHAR(50),
    potronymic VARCHAR(50),
    age SMALLINT,
    gender VARCHAR(10),
    nationality VARCHAR(30),
    created_at TIMESTAMP DEFAULT NOW()
);