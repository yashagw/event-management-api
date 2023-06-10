CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "role" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "password_update_at" timestamp
);

CREATE TABLE "events" (
  "id" bigserial PRIMARY KEY,
  "host_id" bigserial NOT NULL,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "location" varchar NOT NULL,
  "start_date" timestamp NOT NULL,
  "end_date" timestamp NOT NULL,
  "total_tickets" int NOT NULL,
  "left_tickets" int NOT NULL,
  "status" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "tickets" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "event_id" bigserial NOT NULL,
  "quantity" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "user_host_request" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "moderator_id" bigserial NOT NULL,
  "status" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "events" ADD FOREIGN KEY ("host_id") REFERENCES "users" ("id");

ALTER TABLE "user_host_request" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_host_request" ADD FOREIGN KEY ("moderator_id") REFERENCES "users" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");
