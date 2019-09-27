// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"
)

var (
	startConfCh = make(chan int)
	jobDoneCh   = make(chan struct{})
)

func StartJob(confId int) {
	startConfCh <- confId
}

func JobMonitor() {
	var stateJobRunning bool
	for {
		select {
		case confId := <-startConfCh:
			if !stateJobRunning {
				stateJobRunning = true
				go Job(confId)
			}
		case <-jobDoneCh:
			stateJobRunning = false
		}
	}
}

func Job(confId int) {
	defer func() { jobDoneCh <- struct{}{} }()
	if confId <= 0 || confId > len(Cfg.Configs) {
		SetErrorState("Please select a valid configuration!")
	}
	c := Cfg.Configs[confId-1]
	msg := fmt.Sprintf("Start flashing configuration '%s' in flash '%s'...", c.Name, c.FlashChip)
	SetRunningState(msg)
	time.Sleep(2 * time.Second)
	SetSuccessState("Done!")
}
