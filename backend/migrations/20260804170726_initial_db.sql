-- Create "users" table
CREATE TABLE "public"."users" (
  "user_id" text NOT NULL,
  "user_name" character varying(100) NOT NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  PRIMARY KEY ("user_id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email")
);
