-- Create "documents" table
CREATE TABLE "public"."documents" (
  "document_id" text NOT NULL,
  "title" text NOT NULL,
  "content" text NULL,
  "owner_id" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "updated_by" text NOT NULL,
  PRIMARY KEY ("document_id")
);
