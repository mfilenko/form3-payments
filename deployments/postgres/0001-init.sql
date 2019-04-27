DROP TABLE IF EXISTS
  payments,
  payment_types
CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION update_mtime()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE payment_types (
  id SERIAL PRIMARY KEY,
  type character varying(255) UNIQUE NOT NULL CHECK (type <> ''),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payments (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  version integer NOT NULL DEFAULT 0,
  type_id integer NOT NULL REFERENCES payment_types(id),
  organisation_id uuid NOT NULL,
  attributes jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_type_mtime BEFORE UPDATE
  ON payment_types
    FOR EACH row
  EXECUTE PROCEDURE update_mtime();

CREATE TRIGGER update_payment_mtime BEFORE UPDATE
  ON payments
    FOR EACH row
  EXECUTE PROCEDURE update_mtime();
