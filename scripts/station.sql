-- CREATE TABLE "dashboard_snapshot" (
-- 	"id" UUID NOT NULL UNIQUE,
-- 	"name" VARCHAR(255) NOT NULL,
-- 	"org_id" INTEGER NOT NULL,
-- 	"user_id" INTEGER NOT NULL,
-- 	"dashboard" TEXT NOT NULL,
-- 	"expires" TIMESTAMPTZ NOT NULL,
-- 	"created_at" TIMESTAMPTZ NOT NULL,
-- 	"updated_at" TIMESTAMPTZ NOT NULL,
-- 	PRIMARY KEY("id")
-- );
drop database if exists station_management_dev;

-- create user "wqSysadmin" with password 'asd123' inherit;

CREATE DATABASE station_management_dev
    OWNER "wqSysadmin"
    ENCODING 'UTF8'
    LC_COLLATE='en_US.UTF-8'
    LC_CTYPE='en_US.UTF-8'
    TEMPLATE template0;

DROP TABLE IF EXISTS "station" CASCADE;
CREATE TABLE "station" (
	"id" SERIAL NOT NULL UNIQUE,
	"name" VARCHAR(255) NOT NULL,
	"description" TEXT NOT NULL,
	"status" VARCHAR(255) NOT NULL,
	"long" REAL NOT NULL,
	"lat" REAL NOT NULL,
	"country" VARCHAR(50) NOT NULL,
	"station_type" VARCHAR(50),
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"station_manager" INTEGER NOT NULL,
	"water_body_id" INTEGER,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "water_body" CASCADE;
CREATE TABLE "water_body" (
	"id" SERIAL NOT NULL UNIQUE,
	"water_body_type" VARCHAR(255) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"catchment_id" INTEGER NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "parameter" CASCADE;
CREATE TABLE "parameter" (
	"id" SERIAL NOT NULL UNIQUE,
	"name" VARCHAR(255) NOT NULL,
	"unit" VARCHAR(255) NOT NULL,
	"parameter_group" VARCHAR(255) NOT NULL,
	"description" TEXT NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "station_parameter" CASCADE;
CREATE TABLE "station_parameter" (
	"id" SERIAL NOT NULL UNIQUE,
	"parameter_id" INTEGER NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"station_id" INTEGER NOT NULL,
	"last_receiv_at" TIMESTAMPTZ,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "catchment" CASCADE;
CREATE TABLE "catchment" (
	"id" SERIAL NOT NULL UNIQUE,
	"name" VARCHAR(255) NOT NULL,
	"river_basin_district" VARCHAR(255) NOT NULL,
	"country" VARCHAR(255) NOT NULL,
	"updated_at" TIMESTAMPTZ,
	PRIMARY KEY("id")
);


DROP TABLE IF EXISTS "star_dashboard" CASCADE;
CREATE TABLE "star_dashboard" (
	"id" UUID NOT NULL UNIQUE,
	"user_id" INTEGER NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"update_at" TIMESTAMPTZ NOT NULL,
	"version" INTEGER NOT NULL,
	"layout_configuration" TEXT NOT NULL,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "station_key" CASCADE;
CREATE TABLE "station_key" (
	"id" SERIAL NOT NULL UNIQUE,
	"station_id" INTEGER NOT NULL,
	"org_id" INTEGER NOT NULL,
	"is_revoked" SMALLINT NOT NULL DEFAULT 0,
	"name" VARCHAR(255) NOT NULL,
	"key" TEXT NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"update_at" TIMESTAMP NOT NULL,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "station_upload_status" CASCADE;
CREATE TABLE "station_upload_status" (
	"uid" SERIAL NOT NULL UNIQUE,
	"attachment_id" INTEGER NOT NULL,
	"error" TEXT NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"user_id" INTEGER,
	PRIMARY KEY("uid")
);

DROP TABLE IF EXISTS "station_attachments" CASCADE;
CREATE TABLE "station_attachments" (
	"uid" UUID NOT NULL UNIQUE,
	"size" INTEGER NOT NULL,
	"filename" TEXT NOT NULL,
	"content_type" VARCHAR(255) NOT NULL,
	"display_name" TEXT NOT NULL,
	"workflow_state" VARCHAR(255) NOT NULL,
	"user_id" INTEGER NOT NULL,
	"file_state" VARCHAR(255) NOT NULL,
	"namespace" VARCHAR(255) NOT NULL,
	"station_id" INTEGER NOT NULL,
	PRIMARY KEY("uid")
);
ALTER TABLE "station_parameter"
ADD FOREIGN KEY("station_id") REFERENCES "station"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "station_parameter"
ADD FOREIGN KEY("parameter_id") REFERENCES "parameter"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "water_body"
ADD FOREIGN KEY("catchment_id") REFERENCES "catchment"("id")
ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE "station"
ADD FOREIGN KEY("water_body_id") REFERENCES "water_body"("id")
ON UPDATE CASCADE ON DELETE SET NULL;
-- ALTER TABLE "star_dashboard"
-- ADD FOREIGN KEY("user_id") REFERENCES "user"("id")
-- ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE "station_key"
ADD FOREIGN KEY("station_id") REFERENCES "station"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "station_attachments"
ADD FOREIGN KEY("station_id") REFERENCES "station"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
-- ALTER TABLE "station_attachments"
-- ADD FOREIGN KEY("user_id") REFERENCES "user"("id")
-- ON UPDATE NO ACTION ON DELETE NO ACTION;