package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wolftotem4/golava-core/session"
)

type MySQLSessionHandler struct {
	DB    *sql.DB
	Table string
}

func NewMySQLSessionHandler(db *sql.DB, table string) *MySQLSessionHandler {
	return &MySQLSessionHandler{
		DB:    db,
		Table: table,
	}
}

func (d *MySQLSessionHandler) Read(ctx context.Context, sessionId string) ([]byte, error) {
	row := d.DB.QueryRowContext(ctx, fmt.Sprintf(
		"SELECT payload FROM `%s` WHERE id = ?", d.Table,
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

func (d *MySQLSessionHandler) Write(ctx context.Context, sessionId string, data session.SessionData) error {
	now := time.Now().Unix()
	_, err := d.DB.ExecContext(
		ctx,
		fmt.Sprintf(
			"INSERT INTO `%s` (id, user_id, ip_address, user_agent, payload, last_activity) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE user_id = ?, ip_address = ?, user_agent = ?, payload = ?, last_activity = ?",
			d.Table,
		),
		sessionId, data.UserID, data.IPAddress, data.UserAgent, data.Payload, now, data.UserID, data.IPAddress, data.UserAgent, data.Payload, now,
	)
	return err
}

func (d *MySQLSessionHandler) GC(ctx context.Context, lifetime time.Duration) (int64, error) {
	result, err := d.DB.ExecContext(ctx, fmt.Sprintf(
		"DELETE FROM `%s` WHERE last_activity <= ?", d.Table,
	), time.Now().Add(-lifetime).Unix())
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (d *MySQLSessionHandler) Destroy(ctx context.Context, sessionId string) error {
	_, err := d.DB.ExecContext(ctx, fmt.Sprintf(
		"DELETE FROM `%s` WHERE id = ?", d.Table,
	), sessionId)
	return err
}
