-- seed_categories.sql
WITH cat(slug, name, parent_slug, sort_order) AS (
    VALUES
        -- Parents
        ('electronics', 'Electronics', NULL, 1),
        ('books-stationery', 'Books & Stationery', NULL, 2),
        ('fashion', 'Fashion', NULL, 3),
        ('hostel-room-essentials', 'Hostel & Room Essentials', NULL, 4),
        ('sports-fitness', 'Sports & Fitness', NULL, 5),
        ('services', 'Services', NULL, 6),
        ('vehicles-transport', 'Vehicles & Transport', NULL, 7),

        -- Electronics children
        ('phones', 'Phones', 'electronics', 10),
        ('laptops', 'Laptops', 'electronics', 11),
        ('tablets', 'Tablets', 'electronics', 12),
        ('earphones-headphones', 'Earphones / Headphones', 'electronics', 13),
        ('calculators', 'Calculators', 'electronics', 14),
        ('power-banks', 'Power Banks', 'electronics', 15),

        -- Books & Stationery children
        ('textbooks', 'Textbooks', 'books-stationery', 20),
        ('notes-past-questions', 'Notes / Past Questions', 'books-stationery', 21),
        ('stationery', 'Stationery', 'books-stationery', 22),

        -- Fashion children
        ('clothing', 'Clothing', 'fashion', 30),
        ('shoes', 'Shoes', 'fashion', 31),
        ('bags', 'Bags', 'fashion', 32),
        ('accessories', 'Accessories', 'fashion', 33),

        -- Hostel & Room Essentials children
        ('bedding', 'Bedding', 'hostel-room-essentials', 40),
        ('furniture', 'Furniture', 'hostel-room-essentials', 41),
        ('kitchen-items', 'Kitchen Items', 'hostel-room-essentials', 42),
        ('storage-organizers', 'Storage & Organizers', 'hostel-room-essentials', 43),

        -- Sports & Fitness children
        ('gym-equipment', 'Gym Equipment', 'sports-fitness', 50),
        ('sports-gear', 'Sports Gear', 'sports-fitness', 51),

        -- Services children
        ('tutoring', 'Tutoring', 'services', 60),
        ('freelance-skills', 'Freelance / Skills', 'services', 61),
        ('food-catering', 'Food & Catering', 'services', 62),
        ('printing-business-center', 'Printing & Business Center', 'services', 63),

        -- Vehicles & Transport children
        ('vehicles-transport-main', 'Vehicles & Transport', 'vehicles-transport', 70)
),
upserted AS (
    INSERT INTO categories (name, slug, parent_id, sort_order, is_active)
    SELECT c.name, c.slug, p.id, c.sort_order, TRUE
    FROM cat c
    LEFT JOIN categories p ON p.slug = c.parent_slug
    ON CONFLICT (slug) DO UPDATE
      SET name = EXCLUDED.name,
          parent_id = EXCLUDED.parent_id,
          sort_order = EXCLUDED.sort_order,
          is_active = TRUE
    RETURNING id, slug
)
-- Attributes
INSERT INTO category_attributes (category_id, name, label, type, options, required, sort_order)
SELECT c.id, a.name, a.label, a.type, a.options, a.required, a.sort_order
FROM (
    VALUES
        -- Shared product attributes (attach broadly at relevant parents)
        ('electronics', 'condition', 'Condition', 'enum', to_jsonb(ARRAY['new','used','second-hand']), true, 1),
        ('electronics', 'negotiable', 'Negotiable', 'boolean', NULL, false, 2),

        ('books-stationery', 'condition', 'Condition', 'enum', to_jsonb(ARRAY['new','used','second-hand']), true, 1),
        ('books-stationery', 'negotiable', 'Negotiable', 'boolean', NULL, false, 2),

        ('fashion', 'condition', 'Condition', 'enum', to_jsonb(ARRAY['new','used','second-hand']), true, 1),
        ('fashion', 'negotiable', 'Negotiable', 'boolean', NULL, false, 2),

        ('hostel-room-essentials', 'condition', 'Condition', 'enum', to_jsonb(ARRAY['new','used','second-hand']), true, 1),
        ('hostel-room-essentials', 'negotiable', 'Negotiable', 'boolean', NULL, false, 2),

        ('sports-fitness', 'condition', 'Condition', 'enum', to_jsonb(ARRAY['new','used','second-hand']), true, 1),
        ('sports-fitness', 'negotiable', 'Negotiable', 'boolean', NULL, false, 2),

        ('vehicles-transport-main', 'condition', 'Condition', 'enum', to_jsonb(ARRAY['new','used','second-hand']), true, 1),
        ('vehicles-transport-main', 'negotiable', 'Negotiable', 'boolean', NULL, false, 2),

        -- Electronics (parent-level)
        ('electronics', 'brand', 'Brand', 'text', NULL, true, 3),
        ('electronics', 'model', 'Model', 'text', NULL, true, 4),
        ('electronics', 'color', 'Color', 'text', NULL, false, 5),
        ('electronics', 'warranty_remaining', 'Warranty Remaining (months)', 'number', NULL, false, 6),

        -- Phones
        ('phones', 'storage', 'Storage (GB)', 'number', NULL, true, 1),
        ('phones', 'ram', 'RAM (GB)', 'number', NULL, true, 2),
        ('phones', 'battery_health', 'Battery Health (%)', 'number', NULL, false, 3),
        ('phones', 'os', 'OS', 'enum', to_jsonb(ARRAY['Android','iOS','Other']), true, 4),
        ('phones', 'network', 'Network', 'enum', to_jsonb(ARRAY['4G','5G']), true, 5),
        ('phones', 'accessories_included', 'Accessories Included', 'enum', to_jsonb(ARRAY['charger','case','earphones','box']), false, 6),

        -- Laptops
        ('laptops', 'storage', 'Storage (GB)', 'number', NULL, true, 1),
        ('laptops', 'ram', 'RAM (GB)', 'number', NULL, true, 2),
        ('laptops', 'processor', 'Processor', 'text', NULL, true, 3),
        ('laptops', 'os', 'OS', 'enum', to_jsonb(ARRAY['Windows','macOS','Linux','ChromeOS']), true, 4),
        ('laptops', 'screen_size', 'Screen Size (inches)', 'number', NULL, true, 5),
        ('laptops', 'battery_health', 'Battery Health (%)', 'number', NULL, false, 6),
        ('laptops', 'charger_included', 'Charger Included', 'boolean', NULL, true, 7),

        -- Tablets
        ('tablets', 'storage', 'Storage (GB)', 'number', NULL, true, 1),
        ('tablets', 'ram', 'RAM (GB)', 'number', NULL, true, 2),
        ('tablets', 'os', 'OS', 'enum', to_jsonb(ARRAY['Android','iPadOS','Windows','Other']), true, 3),
        ('tablets', 'screen_size', 'Screen Size (inches)', 'number', NULL, true, 4),
        ('tablets', 'charger_included', 'Charger Included', 'boolean', NULL, true, 5),

        -- Earphones / Headphones
        ('earphones-headphones', 'type', 'Type', 'enum', to_jsonb(ARRAY['wired','wireless','TWS']), true, 1),
        ('earphones-headphones', 'noise_cancelling', 'Noise Cancelling', 'boolean', NULL, false, 2),
        ('earphones-headphones', 'battery_life', 'Battery Life (hours)', 'number', NULL, false, 3),

        -- Calculators
        ('calculators', 'type', 'Type', 'enum', to_jsonb(ARRAY['scientific','graphing','basic']), true, 1),
        ('calculators', 'brand', 'Brand', 'text', NULL, true, 2),
        ('calculators', 'model', 'Model', 'text', NULL, true, 3),

        -- Power Banks
        ('power-banks', 'capacity', 'Capacity (mAh)', 'number', NULL, true, 1),
        ('power-banks', 'ports', 'Number of Ports', 'number', NULL, false, 2),
        ('power-banks', 'fast_charging', 'Fast Charging', 'boolean', NULL, false, 3),

        -- Books & Stationery (parent-level)
        ('books-stationery', 'course_relevance', 'Course Relevance', 'text', NULL, false, 3),
        ('books-stationery', 'academic_year', 'Academic Year', 'enum', to_jsonb(ARRAY['100','200','300','400','500','Graduate']), false, 4),
        ('books-stationery', 'institution', 'Institution', 'text', NULL, false, 5),

        -- Textbooks
        ('textbooks', 'title', 'Title', 'text', NULL, true, 1),
        ('textbooks', 'author', 'Author', 'text', NULL, true, 2),
        ('textbooks', 'edition', 'Edition', 'text', NULL, false, 3),
        ('textbooks', 'course_code', 'Course Code', 'text', NULL, true, 4),
        ('textbooks', 'department', 'Department', 'text', NULL, true, 5),
        ('textbooks', 'condition_details', 'Condition Details', 'enum', to_jsonb(ARRAY['highlighted','torn pages','clean']), false, 6),

        -- Notes / Past Questions
        ('notes-past-questions', 'course_code', 'Course Code', 'text', NULL, true, 1),
        ('notes-past-questions', 'department', 'Department', 'text', NULL, true, 2),
        ('notes-past-questions', 'academic_year', 'Academic Year', 'enum', to_jsonb(ARRAY['100','200','300','400','500','Graduate']), true, 3),
        ('notes-past-questions', 'semester', 'Semester', 'enum', to_jsonb(ARRAY['1','2','3']), false, 4),
        ('notes-past-questions', 'format', 'Format', 'enum', to_jsonb(ARRAY['printed','handwritten']), true, 5),

        -- Stationery
        ('stationery', 'item_type', 'Item Type', 'text', NULL, true, 1),

        -- Fashion (parent-level)
        ('fashion', 'size', 'Size', 'text', NULL, true, 3),
        ('fashion', 'color', 'Color', 'text', NULL, false, 4),
        ('fashion', 'gender', 'Gender', 'enum', to_jsonb(ARRAY['male','female','unisex']), false, 5),
        ('fashion', 'material', 'Material', 'text', NULL, false, 6),

        -- Clothing
        ('clothing', 'type', 'Type', 'enum', to_jsonb(ARRAY['tops','bottoms','dresses','suits','sports']), true, 1),
        ('clothing', 'fit', 'Fit', 'enum', to_jsonb(ARRAY['slim','regular','oversized']), false, 2),
        ('clothing', 'occasion', 'Occasion', 'enum', to_jsonb(ARRAY['casual','formal','school','sports']), false, 3),

        -- Shoes
        ('shoes', 'size', 'Size (EU/UK)', 'number', NULL, true, 1),
        ('shoes', 'type', 'Type', 'enum', to_jsonb(ARRAY['sneakers','sandals','formal','boots']), true, 2),
        ('shoes', 'gender', 'Gender', 'enum', to_jsonb(ARRAY['male','female','unisex']), false, 3),

        -- Bags
        ('bags', 'type', 'Type', 'enum', to_jsonb(ARRAY['backpack','tote','handbag','laptop bag']), true, 1),
        ('bags', 'compartments', 'Number of Compartments', 'number', NULL, false, 2),
        ('bags', 'laptop_compartment', 'Laptop Compartment', 'boolean', NULL, false, 3),

        -- Accessories
        ('accessories', 'type', 'Type', 'enum', to_jsonb(ARRAY['watch','belt','jewelry','cap']), true, 1),

        -- Hostel & Room Essentials (parent-level)
        ('hostel-room-essentials', 'brand', 'Brand', 'text', NULL, false, 3),
        ('hostel-room-essentials', 'dimensions', 'Dimensions', 'text', NULL, false, 4),

        -- Bedding
        ('bedding', 'size', 'Size', 'enum', to_jsonb(ARRAY['single','double']), true, 1),
        ('bedding', 'type', 'Type', 'enum', to_jsonb(ARRAY['mattress','sheets','pillow','duvet']), true, 2),
        ('bedding', 'material', 'Material', 'text', NULL, false, 3),

        -- Furniture
        ('furniture', 'type', 'Type', 'enum', to_jsonb(ARRAY['desk','chair','shelf','wardrobe']), true, 1),
        ('furniture', 'material', 'Material', 'text', NULL, false, 2),
        ('furniture', 'assembly_required', 'Assembly Required', 'boolean', NULL, false, 3),

        -- Kitchen Items
        ('kitchen-items', 'type', 'Type', 'enum', to_jsonb(ARRAY['gas','hotplate','kettle','microwave','blender']), true, 1),
        ('kitchen-items', 'brand', 'Brand', 'text', NULL, false, 2),
        ('kitchen-items', 'capacity', 'Capacity', 'text', NULL, false, 3),

        -- Storage & Organizers
        ('storage-organizers', 'type', 'Type', 'enum', to_jsonb(ARRAY['boxes','hangers','drawers','racks']), true, 1),
        ('storage-organizers', 'material', 'Material', 'text', NULL, false, 2),

        -- Sports & Fitness (parent-level)
        ('sports-fitness', 'brand', 'Brand', 'text', NULL, false, 3),

        -- Gym Equipment
        ('gym-equipment', 'type', 'Type', 'enum', to_jsonb(ARRAY['dumbbells','resistance bands','mat','kettlebell']), true, 1),
        ('gym-equipment', 'weight_kg', 'Weight (kg)', 'number', NULL, false, 2),

        -- Sports Gear
        ('sports-gear', 'sport_type', 'Sport Type', 'enum', to_jsonb(ARRAY['football','basketball','tennis','volleyball','running']), true, 1),
        ('sports-gear', 'items_included', 'Items Included', 'text', NULL, false, 2),

        -- Services (parent-level)
        ('services', 'service_type', 'Service Type', 'text', NULL, true, 1),
        ('services', 'availability', 'Availability', 'text', NULL, false, 2),
        ('services', 'delivery_method', 'Delivery Method', 'enum', to_jsonb(ARRAY['online','in-person','both']), true, 3),
        ('services', 'pricing_type', 'Pricing Type', 'enum', to_jsonb(ARRAY['fixed','hourly','negotiable']), true, 4),

        -- Tutoring
        ('tutoring', 'course_code', 'Course Code', 'text', NULL, true, 1),
        ('tutoring', 'department', 'Department', 'text', NULL, true, 2),
        ('tutoring', 'level', 'Level', 'enum', to_jsonb(ARRAY['100','200','300','400','500','Graduate']), true, 3),
        ('tutoring', 'mode', 'Mode', 'enum', to_jsonb(ARRAY['online','in-person']), true, 4),

        -- Freelance / Skills
        ('freelance-skills', 'skill_type', 'Skill Type', 'enum', to_jsonb(ARRAY['graphic design','coding','writing','video editing','photography']), true, 1),
        ('freelance-skills', 'turnaround_time', 'Turnaround Time', 'text', NULL, false, 2),
        ('freelance-skills', 'portfolio_link', 'Portfolio Link', 'text', NULL, false, 3),

        -- Food & Catering
        ('food-catering', 'cuisine_type', 'Cuisine Type', 'text', NULL, true, 1),
        ('food-catering', 'delivery_available', 'Delivery Available', 'boolean', NULL, true, 2),
        ('food-catering', 'order_lead_time', 'Order Lead Time', 'text', NULL, false, 3),
        ('food-catering', 'allergies_info', 'Allergies Info', 'text', NULL, false, 4),

        -- Printing & Business Center
        ('printing-business-center', 'service_type', 'Service Type', 'enum', to_jsonb(ARRAY['printing','binding','lamination','scanning']), true, 1),
        ('printing-business-center', 'price_per_page', 'Price per Page', 'number', NULL, false, 2),

        -- Vehicles & Transport
        ('vehicles-transport-main', 'brand', 'Brand', 'text', NULL, true, 3),
        ('vehicles-transport-main', 'model', 'Model', 'text', NULL, true, 4),
        ('vehicles-transport-main', 'year', 'Year', 'number', NULL, false, 5),
        ('vehicles-transport-main', 'color', 'Color', 'text', NULL, false, 6)
) AS a(slug, name, label, type, options, required, sort_order)
JOIN categories c ON c.slug = a.slug
ON CONFLICT (category_id, name) DO UPDATE
  SET label = EXCLUDED.label,
      type = EXCLUDED.type,
      options = EXCLUDED.options,
      required = EXCLUDED.required,
      sort_order = EXCLUDED.sort_order;