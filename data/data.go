// /home/krylon/go/src/github.com/blicero/dwd/data/data.go
// -*- mode: go; coding: utf-8; -*-
// Created on 24. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-24 12:17:54 krylon>

package data

// Warning represents a weather warning for a specific location and time.
type Warning struct {
	Location      string `json:"regionName"`
	Start         int64  `json:"start"`
	End           int64  `json:"end"`
	Type          int64  `json:"type"`
	State         string `json:"state"`
	Level         int    `json:"level"`
	Description   string `json:"description"`
	Event         string `json:"event"`
	Headline      string `json:"headline"`
	Instruction   string `json:"instruction"`
	StateShort    string `json:"stateShort"`
	AltitudeStart int64  `json:"altitudeStart"`
	AltitudeEnd   int64  `json:"altitudeEnd"`
}

// WeatherInfo represetns an aggregate of warnings issued by the DWD at a given time.
type WeatherInfo struct {
	Time           int64                `json:"time"`
	Warnings       map[string][]Warning `json:"warnings"`
	PrelimWarnings map[string][]Warning `json:"vorabInformation"`
	Copyright      string               `json:"copyright"`
}
