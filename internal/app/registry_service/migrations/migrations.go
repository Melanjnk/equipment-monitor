package migrations

import (
	"github.com/jmoiron/sqlx"	
	_ "github.com/lib/pq"
)

func CreateTableEquipment(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS public.equipment (
			id UUID PRIMARY KEY,
			kind SMALLINT NOT NULL CHECK(kind BETWEEN 0 AND 3),
			status SMALLINT NOT NULL CHECK(status BETWEEN 0 AND 2),
			parameters JSONB NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
			updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp
		);
	`)
	return err
}

func DropTableEquipment(db *sqlx.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS public.equipment`)
	return err
}
