-- CREATE TABLE "dashboard_snapshot" (
-- 	"id" UUID NOT NULL UNIQUE,
-- 	"name" VARCHAR(255) NOT NULL,
-- 	"org_id" INTEGER NOT NULL,
-- 	"user_id" INTEGER NOT NULL,
-- 	"dashboard" TEXT NOT NULL,
-- 	"expires" TIMESTAMP NOT NULL,
-- 	"created_at" TIMESTAMP NOT NULL,
-- 	"updated_at" TIMESTAMP NOT NULL,
-- 	PRIMARY KEY("id")
-- );
-- drop database if exists station_management_dev;
--
-- -- create user "wqSysadmin" with password 'asd123' inherit;
--
-- CREATE DATABASE station_management_dev
--     OWNER "wqSysadmin"
--     ENCODING 'UTF8'
--     LC_COLLATE='en_US.UTF-8'
--     LC_CTYPE='en_US.UTF-8'
--     TEMPLATE template0;

DROP TABLE IF EXISTS "station" CASCADE;
CREATE TABLE "station" (
	"id" SERIAL NOT NULL UNIQUE,
	"name" VARCHAR(255) NOT NULL,
	"description" TEXT NOT NULL,
	"status" VARCHAR(255) NOT NULL,
	"long" REAL NOT NULL,
	"lat" REAL NOT NULL,
	"country" VARCHAR(50),
	"station_type" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"station_manager" INTEGER NOT NULL,
	"water_body_id" INTEGER NOT NULL,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "water_body" CASCADE;
CREATE TABLE "water_body" (
	"id" SERIAL NOT NULL UNIQUE,
	"type" VARCHAR(255) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"catchment_id" INTEGER NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"description" TEXT NOT NULL,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "parameter" CASCADE;
CREATE TABLE "parameter" (
	"id" SERIAL NOT NULL UNIQUE,
	"name" VARCHAR(255) NOT NULL,
	"unit" VARCHAR(255) NOT NULL,
	"parameter_group" VARCHAR(255),
	"description" TEXT,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "station_parameter" CASCADE;
CREATE TABLE "station_parameter" (
	"id" SERIAL NOT NULL UNIQUE,
	"parameter_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"station_id" INTEGER NOT NULL,
	"last_receiv_at" TIMESTAMP,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "catchment" CASCADE;
CREATE TABLE "catchment" (
	"id" SERIAL NOT NULL UNIQUE,
	"name" VARCHAR(255) NOT NULL,
	"river_basin_id" VARCHAR(255) NOT NULL,
	"country" VARCHAR(255) NOT NULL,
    "description" TEXT,
	"updated_at" TIMESTAMP NOT NULL,
	PRIMARY KEY("id")
);
DROP TABLE IF EXISTS "river_basin" CASCADE;
CREATE TABLE "river_basin" (
	"id" SERIAL NOT NULL UNIQUE,
	"name" VARCHAR(255) NOT NULL,
	"description" TEXT NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "country" CASCADE;
CREATE TABLE "country" (
	"id" SERIAL NOT NULL UNIQUE,
	"name" VARCHAR(255) NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	PRIMARY KEY("id")
);

DROP TABLE IF EXISTS "star_dashboard" CASCADE;
CREATE TABLE "star_dashboard" (
	"id" UUID NOT NULL UNIQUE,
	"user_id" INTEGER NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"update_at" TIMESTAMP NOT NULL,
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
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"status" VARCHAR(255) NOT NULL,
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
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
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

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_created_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for table: station
CREATE TRIGGER update_station_updated_at
BEFORE UPDATE OR INSERT ON "station"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Trigger for table: water_body
CREATE TRIGGER update_water_body_updated_at
BEFORE UPDATE OR INSERT ON "water_body"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Trigger for table: parameter
CREATE TRIGGER update_parameter_updated_at
BEFORE UPDATE OR INSERT ON "parameter"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Trigger for table: parameter
CREATE TRIGGER update_parameter_created_at
BEFORE INSERT ON "parameter"
FOR EACH ROW
EXECUTE FUNCTION update_created_at_column();

CREATE TRIGGER update_station_created_at
BEFORE INSERT ON "station"
FOR EACH ROW
EXECUTE FUNCTION update_created_at_column();

-- Trigger for table: station_parameter
CREATE TRIGGER update_station_parameter_updated_at
BEFORE UPDATE OR INSERT ON "station_parameter"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Trigger for table: catchment
CREATE TRIGGER update_catchment_updated_at
BEFORE UPDATE OR INSERT ON "catchment"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Trigger for table: river_basin
CREATE TRIGGER update_river_basin_updated_at
BEFORE UPDATE OR INSERT ON "river_basin"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Trigger for table: country
CREATE TRIGGER update_country_updated_at
BEFORE UPDATE OR INSERT ON "country"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Trigger for table: star_dashboard
CREATE TRIGGER update_star_dashboard_updated_at
BEFORE UPDATE OR INSERT ON "star_dashboard"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Trigger for table: station_key
CREATE TRIGGER update_station_key_updated_at
BEFORE UPDATE OR INSERT ON "station_key"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


select * from river_basin;
select * from catchment;
delete from catchment where id >=0;

select * from water_body;