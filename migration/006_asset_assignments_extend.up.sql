ALTER TABLE asset_assignments
ADD COLUMN IF NOT EXISTS assigned_on TIMESTAMP,
ADD COLUMN IF NOT EXISTS condition_at_assignment TEXT,
ADD COLUMN IF NOT EXISTS notes TEXT,
ADD COLUMN IF NOT EXISTS acknowledgement_status TEXT DEFAULT 'PENDING',
ADD COLUMN IF NOT EXISTS acknowledged_at TIMESTAMP;

UPDATE asset_assignments
SET
    assigned_on = COALESCE(assigned_on, created_at),
    acknowledgement_status = COALESCE(acknowledgement_status, 'PENDING')
WHERE assigned_on IS NULL OR acknowledgement_status IS NULL;
