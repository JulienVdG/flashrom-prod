// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "sync"

// Status describe the current programmer status
type Status string

// Status possible values
const (
	StatusIdle   Status = "idle"
	StatusError  Status = "error"
	StatusSucces Status = "succes"
)

type State struct {
	Status   Status
	Disabled bool
	Config   []string
	ConfigId int    `json:",omitempty"`
	Message  string `json:",omitempty"`
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

func setStateMessage(st Status, message string) {
	stateMu.Lock()
	currentState.Message = message
	currentState.Status = st
	stateMu.Unlock()
	SendCurrentState()
}

func SetErrorState(message string) {
	setStateMessage(StatusError, message)
}

func SetSuccessState(message string) {
	setStateMessage(StatusSucces, message)
}

func UpdateConfigId(confId int) {
	stateMu.Lock()
	currentState.ConfigId = confId
	stateMu.Unlock()
}
