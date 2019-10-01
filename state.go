// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"sync"
	"time"
)

// Status describe the current programmer status
type Status string

// Status possible values
const (
	StatusIdle    Status = "idle"
	StatusError   Status = "error"
	StatusSuccess Status = "success"
	StatusRunning Status = "running"
)

type Log struct {
	Time    time.Time
	Level   Status
	Message string
	Detail  string `json:",omitempty"`
}

type State struct {
	Status   Status
	Disabled bool
	Config   []string
	ConfigId int `json:",omitempty"`
	Message  string
	Logs     []Log
}

var (
	currentState State
	stateMu      sync.RWMutex
)

func GetState() State {
	stateMu.RLock()
	s := currentState
	stateMu.RUnlock()
	return s
}

func setStateMessage(st Status, message, detail string) {
	addLog(st, message, detail)
	stateMu.Lock()
	currentState.Message = message
	currentState.Status = st
	switch st {
	case StatusRunning:
		currentState.Disabled = true
	default:
		currentState.Disabled = false
	}
	stateMu.Unlock()
	SendCurrentState()
}

func SetErrorState(message, detail string) {
	setStateMessage(StatusError, message, detail)
}

func SetSuccessState(message, detail string) {
	setStateMessage(StatusSuccess, message, detail)
}

func SetRunningState(message, detail string) {
	setStateMessage(StatusRunning, message, detail)
}

func UpdateConfigId(confId int) {
	stateMu.Lock()
	currentState.ConfigId = confId
	stateMu.Unlock()
}

func addLog(level Status, message, detail string) {
	log := Log{
		Time:    time.Now(),
		Level:   level,
		Message: message,
		Detail:  detail,
	}
	WriteLog(log)
	stateMu.Lock()
	currentState.Logs = append(currentState.Logs, log)
	stateMu.Unlock()
}

func AddLog(level Status, message, detail string) {
	addLog(level, message, detail)
	SendCurrentState()
}
