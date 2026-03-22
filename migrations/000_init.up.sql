-- ENUMS
CREATE TYPE attendance_status AS ENUM ('PRESENT','ABSENT','LEAVE');
CREATE TYPE leave_status AS ENUM ('PENDING','APPROVED','REJECTED');
CREATE TYPE asset_status AS ENUM ('AVAILABLE','ASSIGNED','UNDER_REPAIR','RETIRED','LOST');

-- AUDIT LOG TABLE
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action TEXT NOT NULL,
    entity TEXT NOT NULL,
    entity_id UUID,
    created_at TIMESTAMP DEFAULT NOW()
);