// /home/krylon/go/src/github.com/blicero/dwd/ui/fyne/gui.go
// -*- mode: go; coding: utf-8; -*-
// Created on 26. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-26 20:03:19 krylon>

// Package fyne implements a GUI using fyne.
package fyne

import "fyne.io/fyne/v2"

// GUI wraps the fyne GUI and its components.
type GUI struct {
	dwdApp    fyne.App
	window    fyne.Window
	layout    fyne.Layout
	container *fyne.Container
}
