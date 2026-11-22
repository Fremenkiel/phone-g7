-- +goose Up
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS events (
  id INTEGER PRIMARY KEY,
  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  direction TEXT NOT NULL,
  status TEXT NOT NULL,
  payload TEXT NOT NULL
  );

CREATE INDEX IF NOT EXISTS queued_events ON events(id)
		WHERE direction = 'out' AND status = 'queued';

-- +goose Down
SELECT 'down SQL query';

DROP TABLE events;
