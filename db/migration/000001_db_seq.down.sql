-- Drop foreign key constraints
ALTER TABLE IF EXISTS "tickets" DROP CONSTRAINT IF EXISTS "tickets_user_id_fkey";
ALTER TABLE IF EXISTS "tickets" DROP CONSTRAINT IF EXISTS "tickets_event_id_fkey";
ALTER TABLE IF EXISTS "user_host_request" DROP CONSTRAINT IF EXISTS "user_host_request_user_id_fkey";
ALTER TABLE IF EXISTS "user_host_request" DROP CONSTRAINT IF EXISTS "user_host_request_moderator_id_fkey";
ALTER TABLE IF EXISTS "events" DROP CONSTRAINT IF EXISTS "events_host_id_fkey";

-- Drop tables
DROP TABLE IF EXISTS "tickets";
DROP TABLE IF EXISTS "user_host_request";
DROP TABLE IF EXISTS "events";
DROP TABLE IF EXISTS "users";