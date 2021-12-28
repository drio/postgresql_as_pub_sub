EXPLAIN ANALYZE
insert into ps_jobs
  (data, status, result, status_change_time)
values
  ((select substr(md5(random()::text), 0, 10)), 'new', -1, current_timestamp);
