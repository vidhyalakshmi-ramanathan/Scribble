-- Modify "documents" table
ALTER TABLE "public"."documents" ALTER COLUMN "created_at" DROP NOT NULL, ALTER COLUMN "updated_at" DROP NOT NULL, ALTER COLUMN "deleted_at" DROP NOT NULL;
-- Create index "idx_documents_deleted_at" to table: "documents"
CREATE INDEX "idx_documents_deleted_at" ON "public"."documents" ("deleted_at");
