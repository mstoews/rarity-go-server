-- Rarity test seed data

-- Categories
INSERT INTO categories (id, name, description, sort_order) VALUES
  ('00000001-0000-0000-0000-000000000001', 'Skincare',    'Serums, moisturisers, SPF and treatments',   1),
  ('00000001-0000-0000-0000-000000000002', 'Makeup',      'Foundation, lip colour, eye and cheek',      2),
  ('00000001-0000-0000-0000-000000000003', 'Fragrance',   'Niche and hard-to-find perfumes',            3),
  ('00000001-0000-0000-0000-000000000004', 'Haircare',    'Scalp treatments, masks and styling',        4),
  ('00000001-0000-0000-0000-000000000005', 'Body',        'Luxury body care and bath rituals',          5)
ON CONFLICT DO NOTHING;

-- Stores
INSERT INTO stores (id, name, address, city, latitude, longitude, website, opening_hours) VALUES
  ('00000002-0000-0000-0000-000000000001', 'Aesop Robson Street',      '1045 Robson St',         'Vancouver',  49.2837, -123.1222, 'https://www.aesop.com', 'Mon–Sat 10–7, Sun 11–6'),
  ('00000002-0000-0000-0000-000000000002', 'Goop Lab Vancouver',       '2208 W 4th Ave',         'Vancouver',  49.2655, -123.1571, 'https://goop.com',      'Mon–Sat 10–6'),
  ('00000002-0000-0000-0000-000000000003', 'Space NK Oakridge',        '650 W 41st Ave',         'Vancouver',  49.2336, -123.1173, 'https://spacenk.com',   'Mon–Sat 10–8, Sun 11–7'),
  ('00000002-0000-0000-0000-000000000004', 'Herriott Grace',           '1-350 E 6th Ave',        'Vancouver',  49.2665, -123.0967, 'https://herriottgrace.com', 'Tue–Sat 11–6'),
  ('00000002-0000-0000-0000-000000000005', 'Nordstrom Beauty Pacific Centre', '799 W Georgia St', 'Vancouver',  49.2825, -123.1196, 'https://nordstrom.com', 'Mon–Sat 10–9, Sun 11–7')
ON CONFLICT DO NOTHING;

-- Cosmetics
INSERT INTO cosmetics (id, name, brand, tagline, description, ingredients, image_url, category_id, is_active) VALUES

  -- Skincare
  ('00000003-0000-0000-0000-000000000001',
   'Parsley Seed Anti-Oxidant Eye Serum', 'Aesop',
   'Brightens and firms the delicate eye area',
   'A lightweight serum that addresses fine lines and puffiness around the eye contour using a potent botanical blend.',
   'Aqua/Water, Glycerin, Camellia Sinensis Leaf Extract, Petroselinum Crispum (Parsley) Seed Extract, Niacinamide',
   NULL,
   '00000001-0000-0000-0000-000000000001', TRUE),

  ('00000003-0000-0000-0000-000000000002',
   'Crème Ancienne Supreme', 'Fresh',
   'The ultimate luxury face cream',
   'Hand-crafted in small batches, this rich cream is inspired by a 900-year-old monastic recipe. Intensely hydrating.',
   'Aqua, Caprylic/Capric Triglyceride, Glycerin, Butyrospermum Parkii (Shea) Butter, Tocopheryl Acetate',
   NULL,
   '00000001-0000-0000-0000-000000000001', TRUE),

  ('00000003-0000-0000-0000-000000000003',
   'Midnight Recovery Concentrate', 'Kiehl''s',
   'Overnight botanical oil for visibly renewed skin',
   'A nightly facial oil formulated with lavender essential oil and squalane to restore skin''s youthful radiance while you sleep.',
   'Squalane, Rosa Canina Fruit Oil, Lavandula Angustifolia (Lavender) Oil, Oenothera Biennis (Evening Primrose) Oil',
   NULL,
   '00000001-0000-0000-0000-000000000001', TRUE),

  -- Makeup
  ('00000003-0000-0000-0000-000000000004',
   'Skin Tint SPF 30', 'RMS Beauty',
   'Sheer, buildable coverage with skincare benefits',
   'A breathable, hydrating skin tint that evens tone while letting your natural texture shine through. Free from synthetic fragrance.',
   'Aqua, Zinc Oxide, Titanium Dioxide, Cocos Nucifera (Coconut) Oil, Rosehip Oil, Aloe Barbadensis',
   NULL,
   '00000001-0000-0000-0000-000000000002', TRUE),

  ('00000003-0000-0000-0000-000000000005',
   'Lip Chic', 'ILIA Beauty',
   'Sheer satin lipstick with skincare complex',
   'A cult-favourite sheer lipstick packed with antioxidants and conditioners. Minimal pigment, maximum shine.',
   'Castor Oil, Caprylic/Capric Triglyceride, Candelilla Wax, Jojoba Esters, Vitamin E',
   NULL,
   '00000001-0000-0000-0000-000000000002', TRUE),

  -- Fragrance
  ('00000003-0000-0000-0000-000000000006',
   'Margiela REPLICA — Jazz Club', 'Maison Margiela',
   'Rum, tobacco leaves and musk',
   'An olfactive memory of a New York jazz bar. Warm, woody and deeply intimate. Part of the iconic REPLICA line.',
   'Alcohol Denat., Aqua, Parfum, Rum, Virginia Cedar, White Musk',
   NULL,
   '00000001-0000-0000-0000-000000000003', TRUE),

  ('00000003-0000-0000-0000-000000000007',
   'Floraiku — One Night in Tokyo', 'Floraiku',
   'Cherry blossom, iris and sandalwood',
   'A rare Japanese niche fragrance with limited global distribution. Inspired by a fleeting evening in Shinjuku.',
   'Alcohol Denat., Parfum, Prunus Serrulata (Cherry Blossom), Iris Pallida, Santalum Album (Sandalwood)',
   NULL,
   '00000001-0000-0000-0000-000000000003', TRUE),

  -- Haircare
  ('00000003-0000-0000-0000-000000000008',
   'Scalp Revival Charcoal + Coconut Oil Micro-Exfoliant Scrub', 'Briogeo',
   'Deep-cleansing scalp detox',
   'A cult scrub that buffs away buildup and flakes without stripping moisture. Hard to find outside specialty retailers.',
   'Aqua, Cocos Nucifera (Coconut) Oil, Charcoal Powder, Mentha Piperita (Peppermint) Oil, Salix Alba (Willow) Bark Extract',
   NULL,
   '00000001-0000-0000-0000-000000000004', TRUE),

  -- Body
  ('00000003-0000-0000-0000-000000000009',
   'Huile Prodigieuse Or', 'Nuxe',
   'Shimmering dry oil for face, body and hair',
   'The legendary French dry oil with a golden shimmer. A pharmacy staple in France, sought after worldwide.',
   'Cyclopentasiloxane, Camellia Sinensis Seed Oil, Rosa Canina Fruit Oil, Mica, Parfum',
   NULL,
   '00000001-0000-0000-0000-000000000005', TRUE),

  ('00000003-0000-0000-0000-000000000010',
   'Scrub de Provençe', 'L''Occitane',
   'Exfoliating shea butter body scrub',
   'A rich, indulgent scrub with lavender from the Valensole plateau and fine sugar granules. Flagship scent.',
   'Sucrose, Prunus Amygdalus Dulcis (Sweet Almond) Oil, Butyrospermum Parkii (Shea) Butter, Lavandula Angustifolia Extract',
   NULL,
   '00000001-0000-0000-0000-000000000005', TRUE)

ON CONFLICT DO NOTHING;

-- Cosmetic ↔ Store availability
INSERT INTO cosmetic_stores (cosmetic_id, store_id, in_stock, notes) VALUES
  -- Aesop eye serum at Aesop Robson
  ('00000003-0000-0000-0000-000000000001', '00000002-0000-0000-0000-000000000001', TRUE,  'Full range available in-store'),
  -- Fresh Crème Ancienne at Space NK & Nordstrom
  ('00000003-0000-0000-0000-000000000002', '00000002-0000-0000-0000-000000000003', TRUE,  NULL),
  ('00000003-0000-0000-0000-000000000002', '00000002-0000-0000-0000-000000000005', TRUE,  NULL),
  -- Kiehl's at Nordstrom & Space NK
  ('00000003-0000-0000-0000-000000000003', '00000002-0000-0000-0000-000000000005', TRUE,  NULL),
  ('00000003-0000-0000-0000-000000000003', '00000002-0000-0000-0000-000000000003', FALSE, 'Call ahead to confirm stock'),
  -- RMS Skin Tint at Goop & Space NK
  ('00000003-0000-0000-0000-000000000004', '00000002-0000-0000-0000-000000000002', TRUE,  NULL),
  ('00000003-0000-0000-0000-000000000004', '00000002-0000-0000-0000-000000000003', TRUE,  NULL),
  -- ILIA Lip Chic at Herriott Grace & Nordstrom
  ('00000003-0000-0000-0000-000000000005', '00000002-0000-0000-0000-000000000004', TRUE,  'Full shade range'),
  ('00000003-0000-0000-0000-000000000005', '00000002-0000-0000-0000-000000000005', TRUE,  NULL),
  -- Margiela Jazz Club at Space NK & Nordstrom
  ('00000003-0000-0000-0000-000000000006', '00000002-0000-0000-0000-000000000003', TRUE,  NULL),
  ('00000003-0000-0000-0000-000000000006', '00000002-0000-0000-0000-000000000005', TRUE,  NULL),
  -- Floraiku at Herriott Grace only
  ('00000003-0000-0000-0000-000000000007', '00000002-0000-0000-0000-000000000004', TRUE,  'Exclusive local stockist — limited units'),
  -- Briogeo scrub at Space NK
  ('00000003-0000-0000-0000-000000000008', '00000002-0000-0000-0000-000000000003', TRUE,  NULL),
  ('00000003-0000-0000-0000-000000000008', '00000002-0000-0000-0000-000000000002', FALSE, 'Out of stock — reorder pending'),
  -- Nuxe dry oil at Space NK & Nordstrom
  ('00000003-0000-0000-0000-000000000009', '00000002-0000-0000-0000-000000000003', TRUE,  NULL),
  ('00000003-0000-0000-0000-000000000009', '00000002-0000-0000-0000-000000000005', TRUE,  NULL),
  -- L'Occitane scrub at Nordstrom
  ('00000003-0000-0000-0000-000000000010', '00000002-0000-0000-0000-000000000005', TRUE,  NULL)

ON CONFLICT DO NOTHING;

-- Test user (password: "password123")
INSERT INTO users (id, email, username, password_hash, sub_status) VALUES
  ('00000004-0000-0000-0000-000000000001',
   'test@rarity.app', 'rarity_tester',
   '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lHHi',
   'free')
ON CONFLICT DO NOTHING;

-- Test subscriber (password: "password123")
INSERT INTO users (id, email, username, password_hash, sub_status, sub_expires_at) VALUES
  ('00000004-0000-0000-0000-000000000002',
   'subscriber@rarity.app', 'rarity_subscriber',
   '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lHHi',
   'active',
   NOW() + INTERVAL '1 year')
ON CONFLICT DO NOTHING;

-- A few reviews from the subscriber
INSERT INTO reviews (cosmetic_id, user_id, rating, text) VALUES
  ('00000003-0000-0000-0000-000000000001', '00000004-0000-0000-0000-000000000002', 5,
   'Genuinely transformative around the eyes. I''ve tried everything and this is the one I keep coming back to.'),
  ('00000003-0000-0000-0000-000000000006', '00000004-0000-0000-0000-000000000002', 5,
   'Jazz Club is the most complimented fragrance I''ve ever worn. Warm, sophisticated, and lasts all evening.'),
  ('00000003-0000-0000-0000-000000000007', '00000004-0000-0000-0000-000000000002', 4,
   'Incredibly rare and beautiful — the cherry blossom note is achingly realistic. Wish it had more longevity.'),
  ('00000003-0000-0000-0000-000000000008', '00000004-0000-0000-0000-000000000002', 5,
   'Cleared my scalp issues in two weeks. The peppermint tingle is incredibly satisfying.')
ON CONFLICT DO NOTHING;
