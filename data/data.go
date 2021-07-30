// /home/krylon/go/src/github.com/blicero/dwd/data/data.go
// -*- mode: go; coding: utf-8; -*-
// Created on 24. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-30 13:51:25 krylon>

package data

import "time"

// Warning represents a weather warning for a specific location and time.
type Warning struct {
	ID            int64
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

// Period returns the timespan the warnings is issued for, as a 2-element array.
// Index 0 is the starting time, index 1 the end.
func (w *Warning) Period() [2]time.Time {
	return [2]time.Time{
		time.Unix(w.Start/1000, 0),
		time.Unix(w.End/1000, 0),
	}
} // func (w *Warning) Period() [2]time.Time

// WarningList is a helper type used for sorting Warnings.
type WarningList []Warning

func (wl WarningList) Len() int      { return len(wl) }
func (wl WarningList) Swap(i, j int) { wl[i], wl[j] = wl[j], wl[i] }
func (wl WarningList) Less(i, j int) bool {
	var w1, w2 *Warning
	w1 = &wl[i]
	w2 = &wl[j]

	if w1.Location == w2.Location {
		return w1.Start < w2.Start
	}

	return w1.Location < w2.Location
} // func (wl WarningList) Less(i, j int) bool

// WeatherInfo represetns an aggregate of warnings issued by the DWD at a given time.
type WeatherInfo struct {
	Time           int64               `json:"time"`
	Warnings       map[int64][]Warning `json:"warnings"`
	PrelimWarnings map[int64][]Warning `json:"vorabInformation"`
	Copyright      string              `json:"copyright"`
}

// TimeStamp returns the time the warnings were last updated.
func (w *WeatherInfo) TimeStamp() time.Time {
	return time.Unix(w.Time/1000, 0)
} // func (w *WeatherInfo) TimeStamp() time.Time
