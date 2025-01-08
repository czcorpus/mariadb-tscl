create table mariadb_tscl_status_monitoring (
  "time" timestamp with time zone NOT NULL,
  instance TEXT,
  threads_connected int,
  max_used_connections int,
  aborted_connects int,
  com_select int,
  com_insert int,
  com_update int,
  com_delete int,
  slow_queries int,
  innodb_buffer_pool_reads int,
  innodb_buffer_pool_read_requests int,
  innodb_row_lock_time int,
  handler_read_first int,
  handler_read_key int,
  handler_read_next int,
  handler_read_rnd int,
  handler_read_rnd_next int,
  bytes_sent int,
  bytes_received int
);
select create_hypertable('mariadb_tscl_status_monitoring', 'time');
