package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateEventsTable, downCreateEventsTable)
}

func upCreateEventsTable(tx *sql.Tx) error {
	query := `
create table if not exists events
(
    id          serial constraint events_pk primary key,
    user_id     int       not null,
    title       varchar   not null,
    description text,
    start_at    timestamp not null,
    end_at      timestamp not null,
    notify_at   timestamp
);
`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func downCreateEventsTable(tx *sql.Tx) error {
	query := `drop table if exists events;`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}
