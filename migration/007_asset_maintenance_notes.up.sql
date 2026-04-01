ALTER TABLE asset_maintenance_logs
ADD COLUMN IF NOT EXISTS notes TEXT;

UPDATE asset_maintenance_logs
SET notes = COALESCE(notes, '')
WHERE notes IS NULL;
