CREATE TABLE gosaas_accounts(
	id SERIAL PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	stripe_id TEXT NOT NULL,
	subscription_id TEXT NOT NULL,
	plan TEXT NOT NULL,
	is_yearly BOOL NOT NULL,
	subscribed_on TIMESTAMP NOT NULL,
	seats INTEGER NOT NULL,
	is_active BOOL NOT NULL
);

CREATE TABLE gosaas_users(
	id SERIAL PRIMARY KEY,
	account_id INTEGER REFERENCES gosaas_accounts(id) ON DELETE CASCADE,
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	token TEXT UNIQUE NOT NULL,
	role INTEGER NOT NULL
);