ALTER TABLE flights
    ADD COLUMN created_by UUID,
    ADD COLUMN last_updated_by UUID,
    ADD COLUMN organization_id UUID;

UPDATE flights
SET
    created_by = COALESCE(created_by, '00000000-0000-0000-0000-000000000002'::uuid),
    last_updated_by = COALESCE(last_updated_by, '00000000-0000-0000-0000-000000000002'::uuid),
    organization_id = COALESCE(organization_id, '00000000-0000-0000-0000-000000000001'::uuid);

ALTER TABLE flights
    ALTER COLUMN created_by SET NOT NULL,
    ALTER COLUMN last_updated_by SET NOT NULL,
    ALTER COLUMN organization_id SET NOT NULL;