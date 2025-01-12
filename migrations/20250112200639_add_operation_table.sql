-- +goose Up
CREATE TABLE operations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  inputs JSON NOT NULL,
  type TEXT NOT NULL,
  result REAL NOT NULL,
  user_id INTEGER NOT NULL,
  created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
  FOREIGN KEY (user_id) REFERENCES users (id),
  CHECK (type IN ('add', 'subtract', 'multiply', 'divide', 'sum'))
);

-- +goose Down
-- +goose StatementBegin
DROP TABLE operations;
-- +goose StatementEnd
