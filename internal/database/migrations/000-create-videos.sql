CREATE TABLE IF NOT EXISTS videos(
  id SERIAL PRIMARY KEY,
  status TEXT NULL,
  external_id varchar(255) UNIQUE NULL,
  url varchar(255) UNIQUE NOT NULL,
  summary TEXT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
