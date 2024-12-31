package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/wolftotem4/golava-core/session"
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

func (d *SqliteSessionHandler) Write(ctx context.Context, sessionId string, data session.SessionData) error {
	now := time.Now().Unix()
	_, err := d.DB.ExecContext(
		ctx,
		"INSERT INTO sessions (id, user_id, ip_address, user_agent, payload, last_activity) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT(id) DO UPDATE SET user_id = $2, ip_address = $3, user_agent = $4, payload = $5, last_activity = $6",
		sessionId, data.UserID, data.IPAddress, data.UserAgent, data.Payload, now,
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
