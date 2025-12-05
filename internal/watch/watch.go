package watch

import (
	"context"
	"database/sql"
	"time"

	"github.com/steipete/imsg/internal/db"
)

// Run polls the database and invokes handler for each new message.
func Run(ctx context.Context, store *sql.DB, chatID int64, startRowID int64, interval time.Duration, handler func(db.Message)) error {
	cursor := startRowID
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		msgs, err := db.MessagesAfter(ctx, store, cursor, chatID, 100)
		if err != nil {
			return err
		}
		for _, m := range msgs {
			handler(m)
			if m.RowID > cursor {
				cursor = m.RowID
			}
		}
		if len(msgs) == 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(interval):
			}
		}
	}
}
