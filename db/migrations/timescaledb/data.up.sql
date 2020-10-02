-- +migrate Up
CREATE TABLE IF NOT EXISTS data_type(
   id smallint NOT NULL PRIMARY KEY,
   name varchar(64) NOT NULL,
   description varchar(2048) NOT NULL,
   created_at timestamp not NULL DEFAULT NOW(),
   updated_at timestamp not NULL DEFAULT NOW()
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS project_data(
   created_at timestamp NOT NULL DEFAULT NOW() PRIMARY KEY,
   project_id integer NOT NULL,
   data_type_id smallint REFERENCES data_type(id),
   run_seq_no integer NOT NULL DEFAULT 0,
   data json
);