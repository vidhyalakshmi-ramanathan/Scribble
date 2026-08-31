-- Create "versions" table
CREATE TABLE "public"."versions" (
  "id" text NOT NULL,
  "document_id" text NOT NULL,
  "version_id" bigint NOT NULL,
  "content" text NOT NULL,
  "modified_by" text NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_versions_document_id" to table: "versions"
CREATE INDEX "idx_versions_document_id" ON "public"."versions" ("document_id");
