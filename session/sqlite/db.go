package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type SqliteSessionHandler struct {
	DB *sql.DB
}

func (d *SqliteSessionHandler) Read(ctx context.Context, sessionId string) ([]byte, error) {
	row := d.DB.QueryRowContext(ctx, "SELECT payload FROM sessions WHERE id = $1", sessionId)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var payload []byte
	err := row.Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return payload, err
}

func (d *SqliteSessionHandler) Write(ctx context.Context, sessionId string, payload []byte) error {
	now := time.Now().Unix()
	_, err := d.DB.ExecContext(
		ctx,
		"INSERT INTO sessions (id, payload, last_activity) VALUES ($1, $2, $3) ON CONFLICT(id) DO UPDATE SET payload = $2, last_activity = $3",
		sessionId, payload, now,
	)
	return err
}

func (d *SqliteSessionHandler) GC(ctx context.Context, lifetime time.Duration) (int64, error) {
	result, err := d.DB.ExecContext(ctx, "DELETE FROM sessions WHERE last_activity <= $1", time.Now().Add(-lifetime).Unix())
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (d *SqliteSessionHandler) Destroy(ctx context.Context, sessionId string) error {
	_, err := d.DB.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", sessionId)
	return err
}
