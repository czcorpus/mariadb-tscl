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
	"time"

	"github.com/czcorpus/hltscl"
)

// Timescalable represents any type which is able
// to export its data in a format required by TimescaleDB writer.
type Timescalable interface {

	// ToTimescaleDB defines a method providing data
	// to be written to a database. The first returned
	// value is for tags, the second one for fields.
	ToTimescaleDB(tableWriter *hltscl.TableWriter) *hltscl.Entry

	// GetTime provides a date and time when the record
	// was created.
	GetTime() time.Time

	// GetTableName provides a destination table name
	GetTableName() string

	// MarshalJSON provides a way how to convert the value into JSON.
	// In APIGuard, this is mostly used for logging and debugging.
	MarshalJSON() ([]byte, error)
}

type ReportingWriter interface {
	LogErrors()
	Write(item Timescalable)
	AddTableWriter(tableName string)
}
