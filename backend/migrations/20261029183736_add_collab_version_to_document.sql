-- Modify "documents" table
ALTER TABLE "public"."documents" ADD COLUMN "collab_version" bigint NOT NULL DEFAULT 0;
-- Rename a column from "version_id" to "version"
ALTER TABLE "public"."versions" RENAME COLUMN "version_id" TO "version";
-- Modify "versions" table
ALTER TABLE "public"."versions" ADD COLUMN "schema_ver" bigint NOT NULL DEFAULT 1;
