CREATE OR REPLACE FUNCTION update_modified_column() RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = current_timestamp;
	RETURN NEW;
END;
$$ language 'plpgsql';
