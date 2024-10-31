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
	"bytes"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type BufferedWriter struct {
	Func   func(line string)
	buffer bytes.Buffer
}

// ---------------------------------------------------------------------------------------
//  public members
// ---------------------------------------------------------------------------------------

func (l *BufferedWriter) Write(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		// todo: make portable
		if p[i] == '\n' {
			l.Flush()
			continue
		}

		l.buffer.Write([]byte{p[i]})
	}

	return len(p), nil
}

func (l *BufferedWriter) Flush() {
	str := l.buffer.String()
	if str != "" {
		l.Func(str)
	}

	l.buffer.Reset()
}
