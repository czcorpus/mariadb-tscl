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

package reporting

import (
	"encoding/json"
	"time"

	"github.com/czcorpus/hltscl"
	"github.com/czcorpus/mariadb-tscl/db"
)

const MariaDBTSCLStatusMonitoringTable = "mariadb_tscl_status_monitoring"

// -----

// BackendActionType represents the most general request type distinction
// independent of a concrete service. Currently we need this mostly
// to monitor actions related to our central authentication, i.e. how
// APIGuard handles unauthenticated users and tries to authenticate them
// (if applicable)
type BackendActionType string

// ----

type ConnectionsStatus struct {
	Created  time.Time `json:"created"`
	Instance string    `json:"instance"`
	db.Status
}

func (status *ConnectionsStatus) ToTimescaleDB(tableWriter *hltscl.TableWriter) *hltscl.Entry {
	return tableWriter.NewEntry(status.Created).
		Str("instance", status.Instance).
		Int("threads_connected", status.ThreadsConnected).
		Int("max_used_connections", status.MaxUsedConnections).
		Int("aborted_connects", status.AbortedConnects).
		Int("com_select", status.ComSelect).
		Int("com_insert", status.ComInsert).
		Int("com_update", status.ComUpdate).
		Int("com_delete", status.ComUpdate).
		Int("slow_queries", status.SlowQueries).
		Int("innodb_buffer_pool_reads", status.InnodbBufferPoolReads).
		Int("innodb_buffer_pool_read_requests", status.InnodbBufferPoolReadRequests).
		Int("innodb_row_lock_time", status.InnodbRowLockTime).
		Int("handler_read_first", status.HandlerReadFirst).
		Int("handler_read_key", status.HandlerReadKey).
		Int("handler_read_next", status.HandlerReadNext).
		Int("handler_read_rnd", status.HandlerReadRnd).
		Int("handler_read_rnd_next", status.HandlerReadRndNext).
		Int("bytes_sent", status.BytesSent).
		Int("bytes_received", status.BytesReceived)
}

func (status *ConnectionsStatus) GetTime() time.Time {
	return status.Created
}

func (status *ConnectionsStatus) GetTableName() string {
	return MariaDBTSCLStatusMonitoringTable
}

func (report *ConnectionsStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(*report)
}
