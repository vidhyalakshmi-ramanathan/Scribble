-- Create "document_invites" table
CREATE TABLE "public"."document_invites" (
  "id" text NOT NULL,
  "document_id" text NOT NULL,
  "invited_email" text NOT NULL,
  "role" text NOT NULL,
  "invited_by" text NOT NULL,
  "status" text NOT NULL,
  "expires_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
