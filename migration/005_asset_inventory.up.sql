
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'asset_type'
    ) THEN
        CREATE TYPE asset_type AS ENUM (
            'LAPTOP','MOBILE','DESKTOP','FURNITURE','OTHER'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'asset_status'
    ) THEN
        CREATE TYPE asset_status AS ENUM (
            'AVAILABLE','ASSIGNED','UNDER_REPAIR','RETIRED','LOST'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'asset_category'
    ) THEN
        CREATE TYPE asset_category AS ENUM (
            'LAPTOP','MOBILE','DESKTOP','CHAIR','TABLE','OTHER'
        );
    END IF;
END$$;
CREATE SEQUENCE IF NOT EXISTS asset_code_seq START 1;

CREATE TABLE IF NOT EXISTS asset_inventory (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  asset_code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  serial_no TEXT,
  asset_type asset_type NOT NULL,
  category asset_category,
  status asset_status DEFAULT 'AVAILABLE',
  brand TEXT,
  model TEXT,
  purchase_date DATE,
  purchase_cost_inr NUMERIC(12,2),
  vendor VARCHAR(200),
  warranty_expiry DATE,
  location TEXT,
  notes TEXT,
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_asset_serial_no
ON asset_inventory(serial_no)
WHERE serial_no IS NOT NULL;


CREATE INDEX IF NOT EXISTS idx_asset_filters
ON asset_inventory(status, asset_type, is_deleted);