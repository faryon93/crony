package util

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
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ---------------------------------------------------------------------------------------
//  constants
// ---------------------------------------------------------------------------------------

var (
	reValidEnvVar = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*=.*$`)
)

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

func LoadEnvFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	env := make([]string, 0)

	lineNo := -1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNo++

		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}

		if !reValidEnvVar.MatchString(scanner.Text()) {
			return nil, fmt.Errorf("malformed environment variable line: %d", lineNo)
		}

		env = append(env, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return env, nil
}
