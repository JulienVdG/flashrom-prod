// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Status describe the current programmer status
type Status string

// Status possible values
const (
	StatusIdle Status = "idle"
)

type State struct {
	Status Status
	Config []string
}

var CurrentState State
