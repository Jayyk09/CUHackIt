-- Reverse seed: remove all pantry items for the first user.
DO $$
DECLARE
    _uid UUID;
BEGIN
    SELECT id INTO _uid FROM users ORDER BY created_at LIMIT 1;
    IF _uid IS NOT NULL THEN
        DELETE FROM pantry_items WHERE user_id = _uid;
    END IF;
END $$;
