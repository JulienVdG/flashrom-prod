// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	//verbDebug = fmt.Printf
	verbDebug = func(string, ...interface{}) {}
	logDebug  = fmt.Printf
	//logDebug = func(string, ...interface{}) {}
)

var (
	logFile io.WriteCloser
	logMu   sync.Mutex
)

func OpenLog(dir string) error {
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return fmt.Errorf("error creating directory '%s': %v", dir, err)
	}
	basename := time.Now().Format(time.RFC3339) + ".log"
	logfilename := filepath.Join(dir, basename)

	logFile, err = os.Create(logfilename)
	if err != nil {
		return fmt.Errorf("error creating file '%s': %v", logfilename, err)
	}

	msg := fmt.Sprintf("Recording logs to %s", logfilename)
	addLog(StatusIdle, "Starting", msg)

	return nil
}

func CloseLog() {
	logMu.Lock()
	defer logMu.Unlock()
	logFile.Close()
}

func WriteLog(line Log) {
	logMu.Lock()
	defer logMu.Unlock()
	msg, err := json.Marshal(&line)
	if err != nil {
		return
	}
	msg = append(msg, '\n')

	logFile.Write([]byte(msg))
}
