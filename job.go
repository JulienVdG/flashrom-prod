// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"regexp"
	"strings"

	expect "github.com/google/goexpect"
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
		msg := fmt.Sprintf("Unexpected error confId=%d is not <= %d", confId, len(Cfg.Configs))
		SetErrorState("Please select a valid configuration!", msg)
	}
	c := Cfg.Configs[confId-1]
	msg := fmt.Sprintf("Start flashing configuration '%s' in flash '%s'...", c.Name, c.FlashChip)
	SetRunningState(msg, c.Commandline)
	cmdline := strings.Fields(c.Commandline)
	gExpect, ch, err := expect.SpawnWithArgs(cmdline, -1, expect.PartialMatch(true))
	if err != nil {
		SetErrorState("Could not run flashrom command!", err.Error())
		return
	}
	defer func() {
		err = gExpect.Close()
		if err != nil {
			logDebug("error closing expect: %v\n", err)
		}
		<-ch
	}()

	stop, i, outlog := jobExpect(gExpect, regexp.MustCompile("(No EEPROM/flash device found.)|(Multiple flash chip definitions match the detected chip\\(s\\):)|(Error: opening file .* failed: No such file or directory)|(Reading old flash chip contents... )"), []string{
		"No Flash chip found! Verify that the chip in the socket is present, correct and correctly oriented.",
		"Multiple chip model are detected, this is a configuration issue!",
		"No flash file found, this is a configuration issue!",
	}, "")
	if stop {
		return
	}
	stop, i, outlog = jobExpect(gExpect, regexp.MustCompile("(FAILED.)|(done.)"), []string{"Failed to read Flash chip! Retry or replace the chip."}, outlog)
	if stop {
		return
	}
	stop, i, outlog = jobExpect(gExpect, regexp.MustCompile("(Erasing and writing flash chip... )"), []string{}, outlog)
	if stop {
		return
	}
	// TODO: do we want to handle the multiple erase function search
	stop, i, outlog = jobExpect(gExpect, regexp.MustCompile("(ERASE FAILED!)|(Warning: Chip content is identical to the requested image.)|(Erase/write done.)"), []string{"Failed to erase Flash chip! Retry or replace the chip."}, outlog)
	if stop {
		return
	}
	verbDebug("i:%v\n", i)
	if i == 2 {
		stop, i, outlog = jobExpect(gExpect, regexp.MustCompile("(Erase/write done.)"), []string{}, outlog)
		if stop {
			return
		}
		SetSuccessState("Done! (Content was identical)", outlog)
		return
	}
	stop, i, outlog = jobExpect(gExpect, regexp.MustCompile("(Verifying flash... )"), []string{}, outlog)
	if stop {
		return
	}
	stop, i, outlog = jobExpect(gExpect, regexp.MustCompile("(Your flash chip is in an unknown state.)|(VERIFIED.)"), []string{"Failed to verify Flash chip! Retry or replace the chip."}, outlog)
	if stop {
		return
	}
	SetSuccessState("Done!", outlog)
}

func jobExpect(gExpect *expect.GExpect, re *regexp.Regexp, messages []string, prevoutlog string) (bool, int, string) {
	out, match, err := gExpect.Expect(re, -1)
	outlog := prevoutlog + out
	verbDebug("Out:%v\nMatch:%#v\nErr:%v\n", out, match, err)
	if err != nil {
		outlog = jobErrorState(gExpect, "flashrom monitoring failed", outlog, err)
		return true, 0, outlog
	}
	for i, m := range messages {
		if match[i+1] != "" {
			outlog = jobErrorState(gExpect, m, outlog, nil)
			return true, i + 1, outlog
		}
	}
	for i := len(messages) + 1; i < len(match); i++ {
		if match[i] != "" {
			return false, i, outlog
		}
	}
	outlog = jobErrorState(gExpect, "flashrom monitoring state failure", outlog, nil)
	return true, -1, outlog
}

func jobErrorState(gExpect *expect.GExpect, message, prevoutlog string, err error) string {
	// Readout output in expect buffer and add them to the log
	out, _, _ := gExpect.Expect(nil, 0)
	outlog := prevoutlog + out
	if err != nil {
		SetErrorState(message, outlog+"\n\n"+err.Error())
		return outlog
	}
	SetErrorState(message, outlog)
	return outlog

}
