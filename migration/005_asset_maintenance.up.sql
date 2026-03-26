-- ENUMS
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'maintenance_type'
    ) THEN
        CREATE TYPE maintenance_type AS ENUM (
            'REPAIR',
            'SERVICE',
            'INSPECTION',
            'DISPOSAL'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'maint_status'
    ) THEN
        CREATE TYPE maint_status AS ENUM (
            'IN_PROGRESS',
            'COMPLETED',
            'SCRAPPED'
        );
    END IF;
END$$;

-- TABLE
CREATE TABLE IF NOT EXISTS asset_maintenance_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    asset_id UUID NOT NULL REFERENCES asset_inventory(id),

    maintenance_type maintenance_type NOT NULL,
    description TEXT NOT NULL,

    sent_for_repair_at DATE,
    vendor VARCHAR(200),

    maint_status maint_status DEFAULT 'IN_PROGRESS',

    returned_from_repair_at DATE,
    repair_cost_inr NUMERIC(12,2),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- INDEX
CREATE INDEX IF NOT EXISTS idx_asset_maintenance 
ON asset_maintenance_logs (asset_id, maint_status, created_at);
