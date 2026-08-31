-- Modify "documents" table
ALTER TABLE "public"."documents" ADD COLUMN "deleted_at" timestamptz NOT NULL;
