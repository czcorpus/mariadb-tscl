// Copyright 2024 Martin Zimandl <martin.zimandl@gmail.com>
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

type Conf struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Status struct {
	ComSelect        int
	ComInsert        int
	ComUpdate        int
	ComDelete        int
	ThreadsConnected int
	SlowQueries      int
	BufferPoolReads  int
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
	rows, err := conn.Query("SHOW GLOBAL STATUS WHERE Variable_name IN ('Com_select','Com_insert','Com_update','Com_delete','Threads_connected','Slow_queries','Innodb_buffer_pool_reads')")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var k string
		var v int
		rows.Scan(&k, &v)
		switch k {
		case "Com_select":
			s.ComSelect = v
		case "Com_insert":
			s.ComInsert = v
		case "Com_update":
			s.ComUpdate = v
		case "Com_delete":
			s.ComDelete = v
		case "Threads_connected":
			s.ThreadsConnected = v
		case "Slow_queries":
			s.SlowQueries = v
		case "Innodb_buffer_pool_reads":
			s.BufferPoolReads = v
		}
	}
	return &s, nil
}
