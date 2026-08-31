-- Create "permissions" table
CREATE TABLE "public"."permissions" (
  "id" text NOT NULL,
  "user_id" text NOT NULL,
  "document_id" text NOT NULL,
  "role" character varying(20) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_document_user" to table: "permissions"
CREATE UNIQUE INDEX "idx_document_user" ON "public"."permissions" ("user_id", "document_id");
