EXPLAIN ANALYZE
insert into ps_jobs
  (data, status, status_change_time)
values
  ((select substr(md5(random()::text), 0, 10)), 'new', current_timestamp);
