// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "time"

type jobState int8

const (
	jobStateIdle jobState = iota
	jobStart
	jobRunning
	jobEnd
)

var (
	jobCh = make(chan jobState)
)

func StartJob() {
	jobCh <- jobStart
}

func JobMonitor() {
	var state jobState
	for {
		select {
		case cmd := <-jobCh:
			switch state {
			case jobStateIdle:
				if cmd == jobStart {
					state = jobRunning
					go Job()
				}
			default:
				if cmd == jobEnd {
					state = jobStateIdle
				}
			}
		}
	}
}

func Job() {
	defer func() { jobCh <- jobEnd }()
	SetRunningState("start")
	time.Sleep(2 * time.Second)
	SetSuccessState("Done!")
}
