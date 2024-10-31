package conf

import (
	"errors"
	"github.com/faryon93/crony/util"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

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

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type Conf struct {
	Jobs []*Job
}

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

// Load loads all cron jobs from the given path.
func Load(path string) (*Conf, error) {
	jobFolders, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	conf := &Conf{
		Jobs: make([]*Job, 0),
	}

	for _, folder := range jobFolders {
		if !folder.IsDir() {
			continue
		}

		jobFolderPath := filepath.Join(path, folder.Name())
		spec, err := loadSpec(jobFolderPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				logrus.Infof("skipping folder '%s': no spec file", folder.Name())
			} else {
				logrus.Errorln("failed to load job group spec file:", err.Error())
			}
			continue
		}

		logrus.Infof("loaded job group '%s': %s", folder.Name(), spec.Cron)
		log := logrus.WithField("folder", folder.Name())

		// find all executable files
		jobs, err := os.ReadDir(jobFolderPath)
		if err != nil {
			log.Errorln("failed to read job group folder:", err.Error())
			continue
		}

		for _, job := range jobs {
			if job.IsDir() {
				continue
			}

			log = log.WithField("job", job.Name())
			fileInfo, err := job.Info()
			if err != nil {
				log.Errorln("failed to read job file info:", err.Error())
				continue
			}

			if !isExecAny(fileInfo.Mode()) {
				continue
			}

			env, err := util.LoadEnvFile(filepath.Join(jobFolderPath, fileInfo.Name()+".env"))
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				log.Errorln("failed to load env file:", err.Error())
			}

			if env != nil {
				log.Infoln("successfully loaded env file")
			}

			conf.Jobs = append(conf.Jobs, &Job{
				Spec: spec,
				Path: filepath.Join(jobFolderPath, job.Name()),
				Env:  env,
			})
		}
	}

	return conf, nil
}

// ---------------------------------------------------------------------------------------
//  private functions
// ---------------------------------------------------------------------------------------

// isExecAny returns true if the given mode has any executable bit set.
func isExecAny(mode os.FileMode) bool {
	return mode&0111 != 0
}
