create table mariadb_tscl_status_monitoring (
  "time" timestamp with time zone NOT NULL,
  instance TEXT,
  com_select int,
  com_insert int,
  com_update int,
  com_delete int,
  threads_connected int,
  slow_queries int,
  innodb_buffer_pool_reads int
);
select create_hypertable('mariadb_tscl_status_monitoring', 'time');
