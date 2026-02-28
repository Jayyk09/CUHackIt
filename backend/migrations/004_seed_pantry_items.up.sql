-- Seed dummy pantry items for development/testing.
-- Uses the FIRST user in the users table so it works regardless of who signed up.

DO $$
DECLARE
    _uid UUID;
BEGIN
    SELECT id INTO _uid FROM users ORDER BY created_at LIMIT 1;

    IF _uid IS NULL THEN
        RAISE NOTICE 'No users found — skipping pantry seed.';
        RETURN;
    END IF;

    -- Delete any existing seed items for this user (idempotent re-run)
    DELETE FROM pantry_items WHERE user_id = _uid;

    INSERT INTO pantry_items
        (user_id, name, brand, category, quantity, unit,
         purchase_date, expiration_date, shelf_life_days,
         calories_per_serving, protein_g, carbs_g, fat_g, fiber_g, sugar_g, sodium_mg, serving_size,
         is_expired, is_expiring_soon)
    VALUES
        -- ====== EXPIRING SOON (within 2 days) – Gemini should prioritise these ======
        (_uid, 'Chicken Breast',      'Perdue',       'MEAT',    2,    'lb',
         CURRENT_DATE - INTERVAL '5 days', CURRENT_DATE + INTERVAL '1 day', 7,
         165, 31, 0, 3.6, 0, 0, 74, '4 oz',
         FALSE, TRUE),

        (_uid, 'Baby Spinach',        'Organic Girl', 'PRODUCE', 5,    'oz',
         CURRENT_DATE - INTERVAL '4 days', CURRENT_DATE + INTERVAL '1 day', 5,
         23, 2.9, 3.6, 0.4, 2.2, 0.4, 79, '3 oz',
         FALSE, TRUE),

        (_uid, 'Greek Yogurt',        'Fage',         'DAIRY',   1,    'item',
         CURRENT_DATE - INTERVAL '10 days', CURRENT_DATE + INTERVAL '2 days', 14,
         100, 17, 6, 0.7, 0, 6, 56, '7 oz',
         FALSE, TRUE),

        (_uid, 'Salmon Fillet',       'Fresh',        'SEAFOOD', 1,    'lb',
         CURRENT_DATE - INTERVAL '2 days', CURRENT_DATE + INTERVAL '1 day', 3,
         208, 20, 0, 13, 0, 0, 59, '4 oz',
         FALSE, TRUE),

        -- ====== FRESH (more than 3 days left) ======
        (_uid, 'Brown Rice',          'Uncle Bens',   'PANTRY',  2,    'lb',
         CURRENT_DATE - INTERVAL '30 days', CURRENT_DATE + INTERVAL '180 days', 365,
         216, 5, 45, 1.8, 3.5, 0.7, 10, '1 cup cooked',
         FALSE, FALSE),

        (_uid, 'Eggs',                'Egglands',     'DAIRY',   12,   'item',
         CURRENT_DATE - INTERVAL '7 days',  CURRENT_DATE + INTERVAL '21 days', 35,
         72, 6, 0.4, 5, 0, 0.2, 71, '1 large',
         FALSE, FALSE),

        (_uid, 'Cheddar Cheese',      'Tillamook',    'DAIRY',   8,    'oz',
         CURRENT_DATE - INTERVAL '14 days', CURRENT_DATE + INTERVAL '30 days', 60,
         110, 7, 0, 9, 0, 0, 180, '1 oz',
         FALSE, FALSE),

        (_uid, 'Pasta (Penne)',       'Barilla',      'PANTRY',  1,    'lb',
         CURRENT_DATE - INTERVAL '60 days', CURRENT_DATE + INTERVAL '300 days', 730,
         200, 7, 42, 1, 2, 2, 0, '2 oz dry',
         FALSE, FALSE),

        (_uid, 'Olive Oil',           'Bertolli',     'PANTRY',  16,   'oz',
         CURRENT_DATE - INTERVAL '90 days', CURRENT_DATE + INTERVAL '270 days', 540,
         120, 0, 0, 14, 0, 0, 0, '1 tbsp',
         FALSE, FALSE),

        (_uid, 'Garlic',              NULL,           'PRODUCE', 1,    'head',
         CURRENT_DATE - INTERVAL '5 days',  CURRENT_DATE + INTERVAL '25 days', 30,
         4, 0.2, 1, 0, 0.1, 0, 1, '1 clove',
         FALSE, FALSE),

        (_uid, 'Tomatoes',            'On the Vine',  'PRODUCE', 4,    'item',
         CURRENT_DATE - INTERVAL '3 days',  CURRENT_DATE + INTERVAL '4 days', 7,
         22, 1, 4.8, 0.2, 1.5, 3.2, 6, '1 medium',
         FALSE, FALSE),

        (_uid, 'Bell Peppers',        NULL,           'PRODUCE', 3,    'item',
         CURRENT_DATE - INTERVAL '4 days',  CURRENT_DATE + INTERVAL '3 days', 7,
         31, 1, 6, 0.3, 2.1, 4.2, 4, '1 medium',
         FALSE, TRUE),

        (_uid, 'Broccoli',            NULL,           'PRODUCE', 1,    'lb',
         CURRENT_DATE - INTERVAL '3 days',  CURRENT_DATE + INTERVAL '4 days', 7,
         55, 3.7, 11, 0.6, 5.1, 2.2, 33, '1 cup',
         FALSE, FALSE),

        (_uid, 'Lemon',               NULL,           'PRODUCE', 2,    'item',
         CURRENT_DATE - INTERVAL '7 days',  CURRENT_DATE + INTERVAL '14 days', 21,
         17, 0.6, 5.4, 0.2, 1.6, 1.5, 1, '1 medium',
         FALSE, FALSE),

        (_uid, 'Soy Sauce',           'Kikkoman',     'PANTRY',  10,   'oz',
         CURRENT_DATE - INTERVAL '60 days', CURRENT_DATE + INTERVAL '300 days', 730,
         10, 1, 1, 0, 0, 0, 920, '1 tbsp',
         FALSE, FALSE);

END $$;
