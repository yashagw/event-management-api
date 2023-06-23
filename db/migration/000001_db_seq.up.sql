CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "role" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "password_updated_at" timestamp NOT NULL DEFAULT('0001-01-01 00:00:00Z')
);

CREATE TABLE IF NOT EXISTS "events" (
  "id" bigserial PRIMARY KEY,
  "host_id" bigint NOT NULL,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "location" varchar NOT NULL,
  "start_date" timestamp NOT NULL,
  "end_date" timestamp NOT NULL,
  "total_tickets" int NOT NULL,
  "left_tickets" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "tickets" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "event_id" bigint NOT NULL,
  "quantity" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "user_host_requests" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "moderator_id" bigint NULL,
  "status" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
  UNIQUE ("user_id")
);

ALTER TABLE "events" ADD FOREIGN KEY ("host_id") REFERENCES "users" ("id");

ALTER TABLE "user_host_requests" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_host_requests" ADD FOREIGN KEY ("moderator_id") REFERENCES "users" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");
