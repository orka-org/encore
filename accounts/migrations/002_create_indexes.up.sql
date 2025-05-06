
CREATE UNIQUE INDEX IF NOT EXISTS accounts_username_idx ON accounts (LOWER(username));
CREATE UNIQUE INDEX IF NOT EXISTS accounts_email_idx ON accounts (LOWER(email));
CREATE UNIQUE INDEX IF NOT EXISTS accounts_phone_idx ON accounts (phone) WHERE phone IS NOT NULL AND phone != '';

-- Index for common queries
CREATE INDEX IF NOT EXISTS accounts_role_idx ON accounts (role);
CREATE INDEX IF NOT EXISTS accounts_created_at_idx ON accounts (created_at);

