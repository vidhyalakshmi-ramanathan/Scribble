-- Modify "documents" table
ALTER TABLE "public"."documents" ALTER COLUMN "version" SET DEFAULT 0, ADD COLUMN "category" text NOT NULL DEFAULT 'General';
-- Modify "versions" table
ALTER TABLE "public"."versions" ADD COLUMN "title" text NOT NULL DEFAULT 'General';
