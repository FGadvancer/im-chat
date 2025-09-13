// Copyright © 2024 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

// OutDir creates the absolute path name from path and checks path exists.
// Returns absolute path including trailing '/' or error if path does not exist.
func OutDir(path string) (string, error) {
	outDir, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	stat, err := os.Stat(outDir)
	if err != nil {
		return "", err
	}

	if !stat.IsDir() {
		return "", fmt.Errorf("output directory %s is not a directory", outDir)
	}
	outDir += "/"
	return outDir, nil
}

func ExitWithError(err error) {
	programName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "%s exit -1: %+v\n\n", programName, err)
	os.Exit(-1)
}

func SIGTERMExit() {
	programName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Warning %s receive process terminal SIGTERM exit 0\n", programName)
}

// ---------- Helpers ----------

func ToMillis(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixMilli()
}

// NewReadableOrderID creates a human-friendly alphanumeric ID.
// Default: "PO-YYYYMMDD-XXXXXX" where X is base36.
// NOTE: replace this with your daily-sequence generator if you adopt the atomic-seq design.
func NewReadableOrderID(_ context.Context) (string, error) {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	day := time.Now().Format("20060102")
	// 6 chars base36 pseudo-random tail
	n := 6
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		v, err := rand.Int(rand.Reader, big.NewInt(36))
		if err != nil {
			return "", err
		}
		buf[i] = alphabet[v.Int64()]
	}
	return fmt.Sprintf("PO-%s-%s", day, string(buf)), nil
}
