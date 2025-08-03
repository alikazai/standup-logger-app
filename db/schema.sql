-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  email VARCHAR(254) NOT NULL UNIQUE,
  password CHAR(60) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Standup entries
CREATE TABLE standup_entries (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  yesterday VARCHAR(2000),
  today VARCHAR(2000),
  blockers VARCHAR(2000),
  date TIMESTAMP DEFAULT NOW()
);
