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

import "github.com/rs/zerolog/log"

type NullWriter struct {
}

func (sw *NullWriter) LogErrors() {
	log.Info().
		Bool("fallbackReporting", true).
		Msg("NullWriter.LogErrors()")
}

func (sw *NullWriter) Write(item Timescalable) {
	log.Info().
		Bool("fallbackReporting", true).
		Any("record", item).
		Msg("NullWriter.Write()")

}

func (sw *NullWriter) AddTableWriter(tableName string) {
	log.Info().
		Bool("fallbackReporting", true).
		Msgf("NullWriter.AddTableWriter(%s)", tableName)
}
