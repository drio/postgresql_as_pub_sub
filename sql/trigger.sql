CREATE OR REPLACE FUNCTION ps_jobs_status_notify()
  RETURNS trigger AS
$$
BEGIN
  PERFORM pg_notify('jobs_status_channel', NEW.id::text);
  /*PERFORM pg_notify('jobs_status_channel', NEW.data::text);*/
  /*PERFORM pg_notify('jobs_status_channel', row_to_json(NEW)::text);*/
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

/* NOTE: It works by polling, use: 
  > LISTEN jobs_status_channel;
  for testing.
 */
CREATE TRIGGER ps_jobs_status
  AFTER INSERT OR UPDATE OF status
  ON ps_jobs
  FOR EACH ROW
EXECUTE PROCEDURE ps_jobs_status_notify();
