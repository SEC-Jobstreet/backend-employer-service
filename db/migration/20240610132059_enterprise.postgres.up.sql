CREATE TABLE "enterprise" (
    "id" uuid PRIMARY KEY,
    "name" varchar,
    "country" varchar,
    "address" varchar,
    "latitude" varchar,
    "longitude" varchar,
    "field" varchar,
    "size" varchar,
    "url" varchar,
    "license" varchar,

    "employer_id" varchar,
    "employer_role" varchar,

    "CreatedAt" timestamptz NOT NULL DEFAULT (now())
);