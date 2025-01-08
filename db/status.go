// Copyright 2024 Martin Zimandl <martin.zimandl@gmail.com>
// Copyright 2025 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2024 Institute of the Czech National Corpus,
//                Faculty of Arts, Charles University
//   This file is part of MARIADB-TSCL.
//
//  MARIADB-TSCL is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  MARIADB-TSCL is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with MARIADB-TSCL.  If not, see <https://www.gnu.org/licenses/>.

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Status struct {
	ThreadsConnected             int `json:"threadsConnected"`
	MaxUsedConnections           int `json:"maxUsedConnections"`
	AbortedConnects              int `json:"abortedConnects"`
	ComSelect                    int `json:"comSelect"`
	ComInsert                    int `json:"comInsert"`
	ComUpdate                    int `json:"comUpdate"`
	ComDelete                    int `json:"comDelete"`
	SlowQueries                  int `json:"slowQueries"`
	InnodbBufferPoolReads        int `json:"innodbBufferPoolReads"`
	InnodbBufferPoolReadRequests int `json:"innodbBufferPoolReadRequests"`
	InnodbRowLockTime            int `json:"innodbRowLockTime"`
	HandlerReadFirst             int `json:"handlerReadFirst"`
	HandlerReadKey               int `json:"handlerReadKey"`
	HandlerReadNext              int `json:"handlerReadNext"`
	HandlerReadRnd               int `json:"handlerReadRnd"`
	HandlerReadRndNext           int `json:"handlerReadRndNext"`
	BytesSent                    int `json:"bytesSent"`
	BytesReceived                int `json:"bytesReceived"`
}

func (conf *Conf) Validate(context string) error {
	if conf.Name == "" && conf.Host == "" && conf.User == "" && conf.Password == "" {
		return errors.New("database not configured")

	} else if conf.Name == "" {
		return fmt.Errorf("%s.name is missing/empty", context)

	} else if conf.Host == "" {
		return fmt.Errorf("%s.host is missing/empty", context)

	} else if conf.User == "" {
		return fmt.Errorf("%s.user is missing/empty", context)

	} else if conf.Password == "" {
		return fmt.Errorf("%s.password is missing/empty", context)
	}
	return nil
}

func OpenDB(conf *Conf) (*sql.DB, error) {
	mconf := mysql.NewConfig()
	mconf.Net = "tcp"
	mconf.Addr = conf.Host
	mconf.User = conf.User
	mconf.Passwd = conf.Password
	mconf.DBName = conf.Name
	mconf.ParseTime = true
	mconf.Loc = time.Local
	mconf.Params = map[string]string{"autocommit": "false"}
	db, err := sql.Open("mysql", mconf.FormatDSN())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetDBStatus(conn *sql.DB) (*Status, error) {
	var s Status
	rows, err := conn.Query(
		"SHOW GLOBAL STATUS WHERE Variable_name IN (" +
			"'Threads_connected', " +
			"'Max_used_connections', " +
			"'Aborted_connects', " + // cummulative
			"'Com_select', " + // cummulative
			"'Com_insert', " + // cummulative
			"'Com_update', " + // cummulative
			"'Com_delete', " + // cummulative
			"'Slow_queries', " + // cummulative
			"'Innodb_buffer_pool_reads', " + // cummulative
			"'Innodb_buffer_pool_read_requests', " + // cummulative
			"'Innodb_row_lock_time', " + // cummulative
			"'Handler_read_first', " + // cummulative
			"'Handler_read_key', " + // cummulative
			"'Handler_read_next', " + // cummulative
			"'Handler_read_rnd', " + // cummulative
			"'Handler_read_rnd_next', " + // cummulative,
			"'Bytes_sent', " + // cummulative
			"'Bytes_received' " + // cummulative
			")")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var k string
		var v int
		rows.Scan(&k, &v)
		switch k {
		case "Threads_connected":
			s.ThreadsConnected = v
		case "Max_used_connections":
			s.MaxUsedConnections = v
		case "Aborted_connects":
			s.AbortedConnects = v
		case "Com_select":
			s.ComSelect = v
		case "Com_insert":
			s.ComInsert = v
		case "Com_update":
			s.ComUpdate = v
		case "Com_delete":
			s.ComDelete = v
		case "Slow_queries":
			s.SlowQueries = v
		case "Innodb_buffer_pool_reads":
			s.InnodbBufferPoolReads = v
		case "Innodb_buffer_pool_read_requests":
			s.InnodbBufferPoolReadRequests = v
		case "Innodb_row_lock_time":
			s.InnodbRowLockTime = v
		case "Handler_read_first":
			s.HandlerReadFirst = v
		case "Handler_read_key":
			s.HandlerReadKey = v
		case "Handler_read_next":
			s.HandlerReadNext = v
		case "Handler_read_rnd":
			s.HandlerReadRnd = v
		case "Handler_read_rnd_next":
			s.HandlerReadRndNext = v
		case "Bytes_sent":
			s.BytesSent = v
		case "Bytes_received":
			s.BytesReceived = v
		}
	}
	return &s, nil
}
