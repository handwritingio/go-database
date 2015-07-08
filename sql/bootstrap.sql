-- migration tracking will live in the meta schema
CREATE SCHEMA IF NOT EXISTS meta;

CREATE TABLE meta.migrations (
  current_version integer PRIMARY KEY NOT NULL,
  filename varchar(1024),
  date_migrated timestamp without time zone NOT NULL DEFAULT now()
);
INSERT INTO meta.migrations (current_version, filename)
  VALUES (0, 'bootstrap.sql');

CREATE TABLE meta.id_types (
  table_name varchar(64) NOT NULL,
  type_id integer NOT NULL,
  PRIMARY KEY (table_name, type_id)
);
CREATE UNIQUE INDEX ON meta.id_types (type_id, table_name);
