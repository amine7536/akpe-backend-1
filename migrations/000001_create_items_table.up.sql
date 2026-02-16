CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO items (name)
SELECT name FROM (VALUES ('alpha'), ('bravo'), ('charlie')) AS v(name)
WHERE NOT EXISTS (SELECT 1 FROM items);
