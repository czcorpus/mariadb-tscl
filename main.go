// Copyright 2024 Martin Zimandl <martin.zimandl@gmail.com>
// Copyright 2024 Institute of the Czech National Corpus,
//                Faculty of Arts, Charles University
//
//  mariadb-tscl is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  mariadb-tscl is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with mariadb-tscl.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/czcorpus/cnc-gokit/logging"
	"github.com/czcorpus/mariadb-tscl/cnf"
	"github.com/czcorpus/mariadb-tscl/db"
	"github.com/czcorpus/mariadb-tscl/general"
	"github.com/rs/zerolog/log"
)

var (
	version   string
	buildDate string
	gitCommit string
)

func main() {
	version := general.VersionInfo{
		Version:   version,
		BuildDate: buildDate,
		GitCommit: gitCommit,
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "MariaDB-TSCL\n\nUsage:\n\t%s [options] start [config.json]\n\t%s [options] version\n",
			filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	action := flag.Arg(0)
	if action == "version" {
		fmt.Printf("mariadb-tscl %s\nbuild date: %s\nlast commit: %s\n", version.Version, version.BuildDate, version.GitCommit)
		return

	} else if action != "start" {
		log.Fatal().Msgf("Unknown action %s", action)
	}
	conf := cnf.LoadConfig(flag.Arg(1))
	logging.SetupLogging(conf.Logging)
	log.Info().Msg("Starting MariaDB-TSCL")

	mariadb, err := db.OpenDB(conf.DB)
	if err != nil {
		log.Fatal().Err(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(10 * time.Second)
	go func(ctx context.Context, conn *sql.DB) {
		for range ticker.C {
			status, err := db.GetDBStatus(conn)
			if err != nil {
				log.Error().Err(err).Send()
			} else {
				log.Info().Any("status", status).Send()
			}
		}
	}(ctx, mariadb)

	<-ctx.Done()
	log.Info().Msg("Stopping...")
	err = mariadb.Close()
	if err != nil {
		log.Error().Err(err).Send()
	}
}
