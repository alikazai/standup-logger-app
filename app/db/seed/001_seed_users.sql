-- db/seed/seed_users.sql

-- Enable pgcrypto for UUID generation if not already done
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Insert demo users
INSERT INTO users (id, name, email, password, created_at) VALUES
  (gen_random_uuid(), 'Alice Johnson', 'alice@example.com', '$2y$12$fakehashedpassword1', NOW()),
  (gen_random_uuid(), 'Bob Smith', 'bob@example.com', '$2y$12$fakehashedpassword2', NOW());
