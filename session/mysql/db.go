package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type MySQLSessionHandler struct {
	DB *sql.DB
}

func (d *MySQLSessionHandler) Read(ctx context.Context, sessionId string) ([]byte, error) {
	row := d.DB.QueryRowContext(ctx, "SELECT payload FROM sessions WHERE id = ?", sessionId)
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

func (d *MySQLSessionHandler) Write(ctx context.Context, sessionId string, payload []byte) error {
	now := time.Now().Unix()
	_, err := d.DB.ExecContext(
		ctx,
		"INSERT INTO sessions (id, payload, last_activity) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE payload = ?, last_activity = ?",
		sessionId, payload, now, payload, now,
	)
	return err
}

func (d *MySQLSessionHandler) GC(ctx context.Context, lifetime time.Duration) (int64, error) {
	result, err := d.DB.ExecContext(ctx, "DELETE FROM sessions WHERE last_activity <= ?", time.Now().Add(-lifetime).Unix())
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (d *MySQLSessionHandler) Destroy(ctx context.Context, sessionId string) error {
	_, err := d.DB.ExecContext(ctx, "DELETE FROM sessions WHERE id = ?", sessionId)
	return err
}
