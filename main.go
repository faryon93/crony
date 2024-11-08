package main

// crony
// Copyright (C) 2024 Maximilian Pachl

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// ---------------------------------------------------------------------------------------
//  imports
// ---------------------------------------------------------------------------------------

import (
	"os"
	"os/exec"
	"syscall"
	"time"

	futil "github.com/faryon93/util"
	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"

	"github.com/faryon93/crony/conf"
	"github.com/faryon93/crony/util"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

func main() {
	logrus.Infoln("starting", GetAppVersion())

	config, err := conf.Load("./")
	if err != nil {
		logrus.Errorln("failed to load configuration:", err.Error())
		os.Exit(1)
	}

	// scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := scheduler.Shutdown()
		if err != nil {
			logrus.Errorln("failed to shutdown scheduler:", err.Error())
			return
		}

		logrus.Infoln("scheduler shutdown successful")
	}()

	for _, job := range config.Jobs {
		logrus.Infof("registering job '%s'", job.Path)

		_, err := scheduler.NewJob(
			gocron.CronJob(job.Spec.Cron, false),
			gocron.NewTask(func() {
				logrus.Infof("starting job '%s'", job.Spec.Name)
				log := logrus.
					WithField("path", job.Path)

				cmd := exec.Command(job.Path)
				cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

				if job.Env != nil {
					cmd.Env = append(os.Environ(), job.Env...)
				}

				if job.DecorateLogs == nil || *job.DecorateLogs {
					stdOutWriter := logWriter(log, "stdout")
					defer stdOutWriter.Flush()
					cmd.Stdout = stdOutWriter

					stdErrWriter := logWriter(log, "stderr")
					defer stdErrWriter.Flush()
					cmd.Stderr = stdErrWriter
				} else {
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
				}

				startTime := time.Now()
				err := cmd.Start()
				if err != nil {
					log.Errorln("failed to start job:", err.Error())
					return
				}

				err = cmd.Wait()
				if err != nil {
					log.Errorln("job failed:", err.Error())
					return
				}

				log.Infof("job '%s' finished in %s", job.Path, time.Since(startTime))
			}),
		)
		if err != nil {
			logrus.WithField("job", job.Path).Errorln("failed to register job:", err.Error())
		}
	}

	// start the scheduler
	scheduler.Start()

	//  wait for shutdown
	futil.WaitSignal(syscall.SIGINT, syscall.SIGTERM)
	logrus.Infoln("received SIGINT/SIGTERM: graceful shutdown")
}

// ---------------------------------------------------------------------------------------
//  private functions
// ---------------------------------------------------------------------------------------

func logWriter(log *logrus.Entry, pipeName string) *util.BufferedWriter {
	return &util.BufferedWriter{
		Func: func(line string) {
			log.WithField("pipe", pipeName).Infoln(line)
		},
	}
}
