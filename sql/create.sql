CREATE TYPE job_status AS ENUM (
  'new',
  'processing',
  'success',
  'error');

CREATE TABLE ps_jobs(
    id SERIAL,
    data varchar(256),
    result int,
    status job_status,
    worker_name varchar(256),
    status_change_time timestamp
);
