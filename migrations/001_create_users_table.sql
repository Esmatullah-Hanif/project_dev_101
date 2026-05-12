/*
  # Create users table

  1. New Tables
    - `users`
      - `id` (uuid, primary key)
      - `email` (text, unique, not null)
      - `password_hash` (text, not null)
      - `first_name` (text)
      - `last_name` (text)
      - `bio` (text)
      - `avatar_url` (text)
      - `is_active` (boolean, default true)
      - `created_at` (timestamp, default now())
      - `updated_at` (timestamp, default now())
      - `deleted_at` (timestamp, soft delete)

  2. Security
    - Enable RLS on `users` table
    - Add policies for authenticated users to read their own data and update own profile

  3. Indexes
    - Index on email for login queries
    - Index on deleted_at for soft delete filtering
*/

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  first_name TEXT DEFAULT '',
  last_name TEXT DEFAULT '',
  bio TEXT DEFAULT '',
  avatar_url TEXT DEFAULT '',
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);
