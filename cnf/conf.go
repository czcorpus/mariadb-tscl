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

package cnf

import (
	"encoding/json"
	"os"
	"time"

	"github.com/czcorpus/cnc-gokit/logging"
	"github.com/czcorpus/mariadb-tscl/db"
	"github.com/czcorpus/mariadb-tscl/reporting"
	"github.com/rs/zerolog/log"
)

// Conf is a global configuration of the app
type Conf struct {
	Logging       logging.LoggingConf `json:"logging"`
	InstanceName  string              `json:"instanceName"`
	CheckInterval time.Duration       `json:"checkInterval"`
	DB            *db.Conf            `json:"db"`
	Reporting     *reporting.Conf     `json:"reporting"`
}

func (conf *Conf) GetLocation() *time.Location { // TODO
	loc, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		log.Fatal().Msg("failed to initialize location")
	}
	return loc
}

func LoadConfig(path string) *Conf {
	if path == "" {
		log.Fatal().Msg("Cannot load config - path not specified")
	}
	rawData, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}
	var conf Conf
	err = json.Unmarshal(rawData, &conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}
	return &conf
}
