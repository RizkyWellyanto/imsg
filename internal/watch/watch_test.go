package watch

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/steipete/imsg/internal/db"
)

func TestRunDeliversMessages(t *testing.T) {
	store := buildDB(t)
	defer func() { _ = store.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	var seen int
	err := Run(ctx, store, 0, 0, 50*time.Millisecond, func(_ db.Message) {
		seen++
	})
	if err != context.DeadlineExceeded {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
	if seen == 0 {
		t.Fatal("expected to see at least one message")
	}
}

func buildDB(t *testing.T) *sql.DB {
	t.Helper()
	dbConn, err := sql.Open("sqlite", "file:watchtest?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	stmts := []string{
		`CREATE TABLE chat (ROWID INTEGER PRIMARY KEY, chat_identifier TEXT, display_name TEXT, service_name TEXT);`,
		`CREATE TABLE message (ROWID INTEGER PRIMARY KEY, handle_id INTEGER, text TEXT, date INTEGER, is_from_me INTEGER, service TEXT);`,
		`CREATE TABLE handle (ROWID INTEGER PRIMARY KEY, id TEXT);`,
		`CREATE TABLE chat_message_join (chat_id INTEGER, message_id INTEGER);`,
		`CREATE TABLE message_attachment_join (message_id INTEGER, attachment_id INTEGER);`,
	}
	for _, s := range stmts {
		if _, err := dbConn.Exec(s); err != nil {
			t.Fatalf("exec %s: %v", s, err)
		}
	}
	now := time.Now().UTC()
	apple := now.Add(-time.Duration(db.AppleEpochOffset) * time.Second).UnixNano()
	_, _ = dbConn.Exec(`INSERT INTO chat(ROWID, chat_identifier) VALUES (1, '+1')`)
	_, _ = dbConn.Exec(`INSERT INTO handle(ROWID, id) VALUES (1, '+1')`)
	_, _ = dbConn.Exec(`INSERT INTO message(ROWID, handle_id, text, date, is_from_me, service) VALUES (1,1,'hi', ?,0,'iMessage')`, apple)
	_, _ = dbConn.Exec(`INSERT INTO chat_message_join(chat_id, message_id) VALUES (1,1)`)
	return dbConn
}
