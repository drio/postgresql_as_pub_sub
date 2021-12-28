UPDATE ps_jobs SET status='running'
WHERE id = (
  SELECT id
  FROM ps_jobs
  WHERE status='new'
  ORDER BY id
  FOR UPDATE SKIP LOCKED
  LIMIT 1
)
RETURNING *; /* Return all the rows from the updated row */
