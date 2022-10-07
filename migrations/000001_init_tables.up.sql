CREATE TABLE IF NOT EXISTS "users" (
"id" bigserial PRIMARY KEY,
"created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
"name" text NOT NULL,
"email" citext UNIQUE NOT NULL,
"password_hash" bytea NOT NULL,
"activated" BOOLEAN NOT NULL,
"version" integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS "posts" (
"id" bigserial PRIMARY KEY,
"created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
"title" text NOT NULL,
"post_text" text NOT NULL,
"img" text NOT NULL,
"read_time" integer NOT NULL,
"liked_by" integer[],
"created_by" bigint NOT NULL REFERENCES users ON DELETE CASCADE,
"version" integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS "comments" (
"id" bigserial PRIMARY KEY,
"created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
"text" text NOT NULL,
"created_by" bigint NOT NULL REFERENCES users ON DELETE CASCADE,
"post_id" bigint NOT NULL REFERENCES posts ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "tokens" (
"hash" bytea PRIMARY KEY,
"user_id" bigint NOT NULL REFERENCES users ON DELETE CASCADE,
"expiry" timestamp(0) with time zone NOT NULL,
"scope" text NOT NULL
);
