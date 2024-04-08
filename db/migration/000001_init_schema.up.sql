CREATE TABLE "applications" (
  "id" bigserial PRIMARY KEY,
  "candidate_id" bigserial NOT NULL,
  "job_id" bigserial NOT NULL,
  "status" varchar NOT NULL,
  "message" varchar NOT NULL DEFAULT '',
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "applications" ("candidate_id");

CREATE INDEX ON "applications" ("job_id");

-- CREATE UNIQUE INDEX ON "applications" ("candidate_id", "job_id");
ALTER TABLE "applications" ADD CONSTRAINT "candidate_job_key" UNIQUE ("candidate_id", "job_id");