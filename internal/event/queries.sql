-- name: ListEvents :many
SELECT id, created_at, direction, status, payload
FROM events ORDER BY id DESC;

-- name: GetPending :many
SELECT id, payload FROM events
WHERE direction = 'out' AND status = 'queued'
ORDER BY id;

-- name: CreateIngoingEvent :exec
INSERT INTO events (direction, status, payload) VALUES ('in', 'received', ?1);

-- name: CreateOutgoingEvent :one
INSERT INTO events (direction, status, payload) VALUES ('out', 'queued', ?1) RETURNING id;

-- name: UpdateSentEvent :exec
UPDATE events SET status = 'sent' WHERE id = ?1;

-- name: RemoveEvents :exec
DELETE FROM events;
