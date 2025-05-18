CREATE TABLE IF NOT EXISTS humans (
    ID SERIAL primary key,
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    potronymic VARCHAR(50),
    age SMALLINT,
    gender VARCHAR(10),
    nationality VARCHAR(30),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_timestamp
BEFORE UPDATE ON humans
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
