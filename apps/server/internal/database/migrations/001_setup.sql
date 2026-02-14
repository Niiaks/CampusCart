-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

-- Create enum types
CREATE TYPE user_role AS ENUM ('user', 'admin');
CREATE TYPE listing_condition AS ENUM ('new', 'used', 'second-hand');
CREATE TYPE media_type AS ENUM ('photo', 'video');
CREATE TYPE feedback_type AS ENUM ('suggestion', 'bug');

---- create above / drop below ----

DROP TYPE IF EXISTS feedback_type;
DROP TYPE IF EXISTS media_type;
DROP TYPE IF EXISTS listing_condition;
DROP TYPE IF EXISTS user_role;
DROP EXTENSION IF EXISTS "citext";
DROP EXTENSION IF EXISTS "uuid-ossp";
