package sqlserver

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type SQLServerSessionHandler struct {
	DB *sql.DB
}

func (d *SQLServerSessionHandler) Read(ctx context.Context, sessionId string) ([]byte, error) {
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

func (d *SQLServerSessionHandler) Write(ctx context.Context, sessionId string, payload []byte) error {
	now := time.Now().Unix()
	_, err := d.DB.ExecContext(
		ctx,
		`
BEGIN tran
	UPDATE sessions WITH (serializable) SET payload = $2, last_activity = $3 WHERE id = $1;
	IF @@rowcount = 0
	BEGIN
		INSERT INTO sessions (id, payload, last_activity) VALUES ($1, $2, $3);
	END
COMMIT tran`,
		sessionId, payload, now,
	)
	return err
}

func (d *SQLServerSessionHandler) GC(ctx context.Context, lifetime time.Duration) (int64, error) {
	result, err := d.DB.ExecContext(ctx, "DELETE FROM sessions WHERE last_activity <= $1", time.Now().Add(-lifetime).Unix())
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (d *SQLServerSessionHandler) Destroy(ctx context.Context, sessionId string) error {
	_, err := d.DB.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", sessionId)
	return err
}
