package conf

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
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// ---------------------------------------------------------------------------------------
//  constants
// ---------------------------------------------------------------------------------------

const (
	specFileName = "crony.hcl"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

// Spec represents the configuration file (crony.hcl) for a group of jobs.
type Spec struct {
	Name string
	Cron string `hcl:"cron"`
}

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

// loadSpec loads the crony Spec from a given path.
func loadSpec(path string) (*Spec, error) {
	spec := Spec{
		Name: filepath.Base(path),
	}

	return &spec, hclsimple.DecodeFile(filepath.Join(path, specFileName), nil, &spec)
}
