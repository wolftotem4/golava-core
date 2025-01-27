package sqlserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wolftotem4/golava-core/session"
)

type SQLServerSessionHandler struct {
	DB    *sql.DB
	Table string
}

func NewSQLServerSessionHandler(db *sql.DB, table string) *SQLServerSessionHandler {
	return &SQLServerSessionHandler{
		DB:    db,
		Table: table,
	}
}

func (d *SQLServerSessionHandler) Read(ctx context.Context, sessionId string) ([]byte, error) {
	row := d.DB.QueryRowContext(ctx, fmt.Sprintf(
		"SELECT payload FROM [%s] WHERE id = @p1", d.Table,
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

func (d *SQLServerSessionHandler) Write(ctx context.Context, sessionId string, data session.SessionData) error {
	now := time.Now().Unix()
	_, err := d.DB.ExecContext(
		ctx,
		fmt.Sprintf(`
BEGIN tran
	UPDATE [%s] WITH (serializable) SET user_id = @p2, ip_address = @p3, user_agent = @p4, payload = @p5, last_activity = @p6 WHERE id = @p1;
	IF @@rowcount = 0
	BEGIN
		INSERT INTO [%[1]s] (id, user_id, ip_address, user_agent, payload, last_activity) VALUES (@p1, @p2, @p3, @p4, @p5, @p6);
	END
COMMIT tran`,
			d.Table,
		),
		sessionId, data.UserID, data.IPAddress, data.UserAgent, data.Payload, now,
	)
	return err
}

func (d *SQLServerSessionHandler) GC(ctx context.Context, lifetime time.Duration) (int64, error) {
	result, err := d.DB.ExecContext(ctx, fmt.Sprintf(
		"DELETE FROM [%s] WHERE last_activity <= @p1", d.Table,
	), time.Now().Add(-lifetime).Unix())
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (d *SQLServerSessionHandler) Destroy(ctx context.Context, sessionId string) error {
	_, err := d.DB.ExecContext(ctx, fmt.Sprintf(
		"DELETE FROM [%s] WHERE id = @p1", d.Table,
	), sessionId)
	return err
}
