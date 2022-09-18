CREATE TABLE IF NOT EXISTS "posts" (
"id" bigserial PRIMARY KEY,
"created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
"title" text NOT NULL,
"post_text" text NOT NULL,
"img" text,
"read_time" int,
"liked_by" int[],
"created_by" int,
"version" integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS "comments" (
"id" bigserial PRIMARY KEY,
"created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
"text" text NOT NULL,
"created_by" int,
"post_id" int,
"version" integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS "users" (
"id" bigserial PRIMARY KEY,
"created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
"name" text NOT NULL,
"email" citext UNIQUE NOT NULL,
"password_hash" bytea NOT NULL,
"activated" BOOLEAN NOT NULL,
"version" integer NOT NULL DEFAULT 1
);