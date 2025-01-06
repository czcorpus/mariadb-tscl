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
	"context"
	"fmt"
	"time"

	"github.com/czcorpus/hltscl"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Table struct {
	writer    *hltscl.TableWriter
	opsDataCh chan<- hltscl.Entry
	errCh     <-chan hltscl.WriteError
}

type TimescaleDBWriter struct {
	ctx    context.Context
	tz     *time.Location
	conn   *pgxpool.Pool
	tables map[string]*Table
}

func (sw *TimescaleDBWriter) LogErrors() {
	for name, table := range sw.tables {
		go func(name string, table *Table) {
			for {
				select {
				case <-sw.ctx.Done():
					log.Info().Msgf("about to close %s status writer", name)
					return
				case err := <-table.errCh:
					log.Error().
						Err(err.Err).
						Str("entry", err.Entry.String()).
						Msg("error writing data to TimescaleDB")
					fmt.Println("reporting timescale write err: ", err.Err)
				}
			}
		}(name, table)
	}
}

func (sw *TimescaleDBWriter) Write(item Timescalable) {
	table, ok := sw.tables[item.GetTableName()]
	if ok {
		table.opsDataCh <- *item.ToTimescaleDB(table.writer)
	} else {
		log.Warn().Str("table_name", item.GetTableName()).Msg("Undefined table name in writer")
	}
}

func (sw *TimescaleDBWriter) AddTableWriter(tableName string) {
	twriter := hltscl.NewTableWriter(sw.conn, tableName, "time", sw.tz)
	opsDataCh, errCh := twriter.Activate()
	sw.tables[tableName] = &Table{
		writer:    twriter,
		opsDataCh: opsDataCh,
		errCh:     errCh,
	}
}

func NewReportingWriter(connection *pgxpool.Pool, tz *time.Location, ctx context.Context) *TimescaleDBWriter {
	return &TimescaleDBWriter{
		ctx:    ctx,
		tz:     tz,
		conn:   connection,
		tables: make(map[string]*Table),
	}
}
