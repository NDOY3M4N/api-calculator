-- +goose Up
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pseudo TEXT NOT NULL
);

INSERT INTO users (pseudo) VALUES
  ('p4p1'),
  ('b4tm4n');

-- +goose Down
DELETE FROM users;

DROP TABLE users;
