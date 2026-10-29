-- Modify "documents" table
ALTER TABLE "public"."documents" ALTER COLUMN "content" TYPE jsonb USING content::jsonb, ALTER COLUMN "content" SET NOT NULL, ADD COLUMN "search_text" text NULL, ADD COLUMN "schema_ver" bigint NOT NULL DEFAULT 1;

-- Modify "versions" table
ALTER TABLE "public"."versions" ALTER COLUMN "content" TYPE jsonb USING content::jsonb, ALTER COLUMN "title" DROP DEFAULT;