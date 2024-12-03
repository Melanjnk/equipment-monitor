package migrations

import (
	"github.com/jmoiron/sqlx"	
	_ "github.com/lib/pq"
)

func CreateAutoUpdate(db *sqlx.DB) error {
	_, err := db.Exec(
`CREATE OR REPLACE FUNCTION update_modified_column() RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = current_timestamp;
	RETURN NEW;
END;
$$ language 'plpgsql';`)
	return err
}

func CreateGeneratorUUIDv6(db *sqlx.DB) error {
	// Source: https://gist.github.com/fabiolimace/515a0440e3e40efeb234e12644a6a346
	_, err := db.Exec(
`CREATE OR REPLACE FUNCTION gen_random_uuid_v6() RETURNS uuid AS $$
DECLARE
BEGIN
	RETURN uuid6(clock_timestamp());
END $$ language plpgsql;

CREATE OR REPLACE FUNCTION uuid6(p_timestamp timestamp WITH time zone) RETURNS uuid AS $$
DECLARE
	v_time double precision := NULL;

	v_gregorian_t bigint := NULL;
	v_clock_sequence_and_node bigint := NULL;

	v_gregorian_t_hex_a varchar := NULL;
	v_gregorian_t_hex_b varchar := NULL;
	v_clock_sequence_and_node_hex varchar := NULL;

	c_epoch double precision := 12219292800; -- RFC-9562 epoch: 1582-10-15
	c_100ns_factor double precision := 10^7; -- RFC-9562 precision: 100 ns

	c_version bigint := x'0000000000006000'::bigint; -- RFC-9562 version: b'0110...'
	c_variant bigint := x'8000000000000000'::bigint; -- RFC-9562 variant: b'10xx...'
BEGIN
	v_time := EXTRACT(epoch FROM p_timestamp);

	v_gregorian_t := TRUNC((v_time + c_epoch) * c_100ns_factor);
	v_clock_sequence_and_node := TRUNC(RANDOM() * 2^30)::bigint << 32 | TRUNC(RANDOM() * 2^32)::bigint;

	v_gregorian_t_hex_a := LPAD(to_hex((v_gregorian_t >> 12)), 12, '0');
	v_gregorian_t_hex_b := LPAD(to_hex((v_gregorian_t & 4095) | c_version), 4, '0');
	v_clock_sequence_and_node_hex := LPAD(to_hex(v_clock_sequence_and_node | c_variant), 16, '0');

	RETURN (v_gregorian_t_hex_a || v_gregorian_t_hex_b  || v_clock_sequence_and_node_hex)::uuid;
END $$ language plpgsql;`)
	return err
}

func CreateTableEquipment(db *sqlx.DB) error {
	_, err := db.Exec(
`CREATE TABLE IF NOT EXISTS public.equipment (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid_v6(),
	kind SMALLINT NOT NULL CHECK(kind BETWEEN 0 AND 3),
	status SMALLINT NOT NULL CHECK(status BETWEEN 0 AND 2) DEFAULT 0,
	parameters JSONB NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp CHECK (updated_at >= created_at)
);

CREATE TRIGGER update_modified_time BEFORE UPDATE ON public.equipment FOR EACH ROW EXECUTE PROCEDURE update_modified_column();`)
	return err
}

func DropTableEquipment(db *sqlx.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS public.equipment`)
	return err
}
