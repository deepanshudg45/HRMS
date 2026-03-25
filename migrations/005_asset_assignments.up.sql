CREATE TABLE IF NOT EXISTS asset_assignments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  asset_id UUID REFERENCES asset_inventory(id),
  assigned_to UUID,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW()
);

-- PARTIAL UNIQUE INDEX
CREATE UNIQUE INDEX IF NOT EXISTS ux_active_assignment
ON asset_assignments(asset_id)
WHERE is_active = TRUE;