 CREATE TABLE files (
     id SERIAL PRIMARY KEY,
     file_name varchar(255) unique,
     file_data bytea NOT NULL,
     uploaded_at bigint,
     updated_at bigint
 )