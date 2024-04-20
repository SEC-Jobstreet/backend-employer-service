CREATE TABLE employer_profile (
  "id" bigserial PRIMARY KEY,
  "email" varchar(255) NOT NULL UNIQUE,
  "first_name" varchar(255) NOT NULL,
  "last_name" varchar(255) NOT NULL,
  "phone" varchar(12) NOT NULL,
  "address" varchar(255) NOT NULL,
  "email_confirmed" boolean NOT NULL DEFAULT false,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE enterprise (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "country" varchar(255),
  -- Mock location data
  "location_lat" FLOAT,
  "location_lon" FLOAT,
  "field" varchar(255),
  "size" varchar(255),
  "role" varchar(255),
  "url" varchar(255),
  "gpkd" varchar(255),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE employer_enterprise (
  employer_id bigint NOT NULL,
  enterprise_id bigint NOT NULL,
  PRIMARY KEY (employer_id, enterprise_id),
  FOREIGN KEY (employer_id) REFERENCES employer_profile(id),
  FOREIGN KEY (enterprise_id) REFERENCES enterprise(id)
);

CREATE INDEX ON "employer_enterprise" ("employer_id");
CREATE INDEX ON "employer_enterprise" ("enterprise_id");
