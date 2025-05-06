CREATE TABLE IF NOT EXISTS accounts(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	
	username TEXT NOT NULL,
	email TEXT NOT NULL,
	password TEXT NOT NULL,
	
	phone TEXT,
	
	first_name TEXT,
	last_name TEXT,
	
	role TEXT NOT NULL CHECK (role IN ('user', 'admin', 'moderator')),
	
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

