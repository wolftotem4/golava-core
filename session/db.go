package session

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type DatabaseSessionHandler struct {
	DB         *sql.DB
	DriverName string
}

func (d *DatabaseSessionHandler) Read(ctx context.Context, sessionId string) ([]byte, error) {
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

func (d *DatabaseSessionHandler) Write(ctx context.Context, sessionId string, payload []byte) error {
	switch d.DriverName {
	case "sqlite", "sqlite3":
		_, err := performSqliteUpsert(ctx, d.DB, sessionId, payload)
		return err
	case "mysql":
		_, err := performMysqlUpsert(ctx, d.DB, sessionId, payload)
		return err
	default:
		_, err := performGenericUpsert(ctx, d.DB, sessionId, payload)
		return err
	}
}

func (d *DatabaseSessionHandler) GC(ctx context.Context, lifetime time.Duration) (int64, error) {
	result, err := d.DB.ExecContext(ctx, "DELETE FROM sessions WHERE last_activity <= ?", time.Now().Add(-lifetime).Unix())
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (d *DatabaseSessionHandler) Destroy(ctx context.Context, sessionId string) error {
	_, err := d.DB.ExecContext(ctx, "DELETE FROM sessions WHERE id = ?", sessionId)
	return err
}

func performSqliteUpsert(ctx context.Context, db *sql.DB, sessionId string, payload []byte) (sql.Result, error) {
	now := time.Now().Unix()
	return db.ExecContext(ctx, "INSERT INTO sessions (id, payload, last_activity) VALUES (?, ?, ?) ON CONFLICT(id) DO UPDATE SET payload = ?, last_activity = ?", sessionId, payload, now, payload, now)
}

func performMysqlUpsert(ctx context.Context, db *sql.DB, sessionId string, payload []byte) (sql.Result, error) {
	now := time.Now().Unix()
	return db.ExecContext(ctx, "INSERT INTO sessions (id, payload, last_activity) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE payload = ?, last_activity = ?", sessionId, payload, now, payload, now)
}

func performGenericUpsert(ctx context.Context, db *sql.DB, sessionId string, payload []byte) (sql.Result, error) {
	if exists, err := sessionIdExists(ctx, db, sessionId); err != nil {
		return nil, err
	} else if exists {
		return performUpdate(ctx, db, sessionId, payload)
	}
	return performInsert(ctx, db, sessionId, payload)
}

func sessionIdExists(ctx context.Context, db *sql.DB, sessionId string) (bool, error) {
	row := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions WHERE id = ?", sessionId)
	if err := row.Err(); err != nil {
		return false, err
	}

	var count int
	err := row.Scan(&count)
	return count > 0, err
}

func performUpdate(ctx context.Context, db *sql.DB, sessionId string, payload []byte) (sql.Result, error) {
	now := time.Now().Unix()
	return db.ExecContext(ctx, "UPDATE sessions SET payload = ?, last_activity = ? WHERE id = ?", payload, now, sessionId)
}

func performInsert(ctx context.Context, db *sql.DB, sessionId string, payload []byte) (sql.Result, error) {
	now := time.Now().Unix()
	return db.ExecContext(ctx, "INSERT INTO sessions (id, payload, last_activity) VALUES (?, ?, ?)", sessionId, payload, now)
}
