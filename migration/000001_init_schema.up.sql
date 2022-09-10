CREATE TABLE IF NOT EXISTS "posts" (
"id" bigserial PRIMARY KEY,
"created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
"title" text NOT NULL,
"post_text" text NOT NULL,
"img" text,
"read_time" int,
"liked_by" int[],
"created_by" int
);

CREATE TABLE IF NOT EXISTS "comments" (
"id" bigserial PRIMARY KEY,
"created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
"text" text NOT NULL,
"created_by" int,
"post_id" int
);

CREATE TABLE IF NOT EXISTS "users" (
"id" bigserial PRIMARY KEY,
"created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
"name" text NOT NULL,
"email" text NOT NULL,
"password" text NOT NULL
);