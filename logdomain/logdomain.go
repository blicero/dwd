// /home/krylon/go/src/github.com/blicero/dwd/logdomain/logdomain.go
// -*- mode: go; coding: utf-8; -*-
// Created on 23. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-23 17:25:23 krylon>

// Package logdomain provides constants to identify the different
// "areas" of the application that perform logging.
package logdomain

//go:generate stringer -type=ID

// ID represents an area of concern.
type ID uint8

// These constants identify the various logging domains.
const (
	Common ID = iota
	Client
	DBPool
	Database
	Web
)

// AllDomains returns a slice of all the known log sources.
func AllDomains() []ID {
	return []ID{
		Common,
		Client,
		DBPool,
		Database,
		Web,
	}
} // func AllDomains() []ID
