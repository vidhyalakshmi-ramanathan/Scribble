-- Modify "documents" table
ALTER TABLE "public"."documents" ADD COLUMN "version" bigint NOT NULL DEFAULT 1;
-- Drop index "idx_versions_document_id" from table: "versions"
DROP INDEX "public"."idx_versions_document_id";
-- Create index "idx_document_version" to table: "versions"
CREATE UNIQUE INDEX "idx_document_version" ON "public"."versions" ("document_id", "version_id");
