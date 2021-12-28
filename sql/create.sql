CREATE TYPE job_status AS ENUM (
  'new',
  'running',
  'success',
  'error');

CREATE TABLE ps_jobs(
    id SERIAL,
    data varchar(256),
    status job_status,
    status_change_time timestamp
);
