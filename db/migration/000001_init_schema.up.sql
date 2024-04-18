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

CREATE TABLE employer_profile (
  "id" bigserial PRIMARY KEY,
  "enterprise_id" bigserial NOT NULL,
  "email" varchar(255) NOT NULL UNIQUE,
  "first_name" varchar(255) NOT NULL,
  "last_name" varchar(255) NOT NULL,
  "phone" varchar(255),
  "address" TEXT,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "employer_profile" ("enterprise_id");

CREATE TABLE enterprise (
  "id" bigserial PRIMARY KEY,
  "employer_id" bigserial UNIQUE,
  "name" varchar(255) NOT NULL,
  "country" varchar(255),
  -- Assuming location (lat, lon) should be separate columns
  "location_lat" FLOAT,
  "location_lon" FLOAT,
  "field" varchar(255),
  "size" varchar(255),
  "role" varchar(255),
  "url" varchar(255),
  "gpkd" varchar(255), -- Likely a custom code or identifier
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "enterprise" ("employer_id");
