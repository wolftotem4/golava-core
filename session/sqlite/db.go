package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wolftotem4/golava-core/session"
)

type SqliteSessionHandler struct {
	DB    *sql.DB
	Table string
}

func NewSqliteSessionHandler(db *sql.DB, table string) *SqliteSessionHandler {
	return &SqliteSessionHandler{
		DB:    db,
		Table: table,
	}
}

func (d *SqliteSessionHandler) Read(ctx context.Context, sessionId string) ([]byte, error) {
	row := d.DB.QueryRowContext(ctx, fmt.Sprintf(
		`SELECT payload FROM "%s" WHERE id = $1`, d.Table,
	), sessionId)
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
		fmt.Sprintf(
			`INSERT INTO "%s" (id, user_id, ip_address, user_agent, payload, last_activity) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT(id) DO UPDATE SET user_id = $2, ip_address = $3, user_agent = $4, payload = $5, last_activity = $6`,
			d.Table,
		),
		sessionId, data.UserID, data.IPAddress, data.UserAgent, data.Payload, now,
	)
	return err
}

func (d *SqliteSessionHandler) GC(ctx context.Context, lifetime time.Duration) (int64, error) {
	result, err := d.DB.ExecContext(ctx, fmt.Sprintf(
		`DELETE FROM "%s" WHERE last_activity <= $1`,
		d.Table,
	), time.Now().Add(-lifetime).Unix())
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (d *SqliteSessionHandler) Destroy(ctx context.Context, sessionId string) error {
	_, err := d.DB.ExecContext(ctx, fmt.Sprintf(
		`DELETE FROM "%s" WHERE id = $1`,
		d.Table,
	), sessionId)
	return err
}
