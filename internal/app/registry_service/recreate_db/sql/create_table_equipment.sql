CREATE TABLE IF NOT EXISTS public.equipment (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid_v6(),
	kind SMALLINT NOT NULL CHECK(kind BETWEEN 0 AND 3),
	status SMALLINT NOT NULL CHECK(status BETWEEN 0 AND 2) DEFAULT 0,
	parameters JSONB NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp CHECK (updated_at >= created_at)
);

CREATE TRIGGER update_modified_time BEFORE UPDATE ON public.equipment FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
