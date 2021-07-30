// /home/krylon/go/src/github.com/blicero/dwd/ui/fyne/gui.go
// -*- mode: go; coding: utf-8; -*-
// Created on 26. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-29 22:16:35 krylon>

// +build !386

// Package fyne implements a GUI using fyne.
package fyne

import (
	"fmt"
	"log"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/blicero/dwd/common"
	"github.com/blicero/dwd/data"
	"github.com/blicero/dwd/logdomain"
)

// GUI wraps the fyne GUI and its components.
type GUI struct {
	lock          sync.RWMutex
	warnMap       map[int64]bool
	warnings      []data.Warning
	log           *log.Logger
	client        *data.Client
	dwdApp        fyne.App
	window        fyne.Window
	container     *fyne.Container
	refreshButton *widget.Button
	lst           *widget.List
	//tbl           *widget.Table
}

func MakeGUI(c *data.Client) (*GUI, error) {
	var (
		err error
		g   = &GUI{
			client: c,
		}
	)

	if g.log, err = common.GetLogger(logdomain.GUI); err != nil {
		return nil, err
	} /* else if g.client, err = data.New("", patterns...); err != nil {
		g.log.Printf("[ERROR] Cannot create Client: %s\n",
			err.Error())
		return nil, err
	} */

	// defer g.client.Start()

	g.dwdApp = app.NewWithID("dwd_notifier")
	g.warnMap = make(map[int64]bool)
	g.refreshButton = widget.NewButton("Aktualisieren", g.refresh)

	// g.tbl = widget.NewTable(
	// 	g.tblSize,
	// 	func() fyne.CanvasObject { return widget.NewLabel("") },
	// 	g.tblUpdate,
	// )

	g.lst = widget.NewList(
		g.wcnt,
		func() fyne.CanvasObject { return widget.NewLabel("") },
		g.lstUpdate)

	g.container = container.NewVBox(
		g.refreshButton,
		g.lst,
	)

	g.window = g.dwdApp.NewWindow("DWD")
	g.window.SetContent(g.container)

	return g, nil
} // func MakeGUI(c *data.Client) (*GUI, error)

func (g *GUI) ShowAndRun() {
	go g.wloop()
	g.window.ShowAndRun()
} // func (g *GUI) ShowAndRun()

func (g *GUI) wcnt() int {
	g.lock.RLock()
	var cnt = len(g.warnings)
	g.lock.RUnlock()
	return cnt
} // func (g *GUI) wcnt() int

// func (g *GUI) tblSize() (int, int) {
// 	g.lock.RLock()
// 	defer g.lock.RUnlock()
// 	var rows = len(g.warnings)
// 	const cols = 4

// 	g.log.Printf("[DEBUG] tblSize -> (%d rows, %d cols)\n",
// 		rows,
// 		cols)
// 	return cols, rows
// } // func (g *GUI) tblSize() (int, int)

func (g *GUI) lstUpdate(id widget.ListItemID, obj fyne.CanvasObject) {
	// g.log.Printf("[TRACE] tblUpdate -> Row = %d, Col = %d\n",
	// 	id.Row,
	// 	id.Col)

	// var idx = id.Row

	var idx = int(id)

	if idx < 0 || idx >= len(g.warnings) {
		g.log.Printf("[ERROR] Invalid Row ID %d, valid range is 0 - %d\n",
			idx,
			len(g.warnings)-1)
		return
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	var (
		lbl *widget.Label
		val string
		w   = g.warnings[idx]
	)

	switch l := obj.(type) {
	case *widget.Label:
		lbl = l
	default:
		g.log.Printf("[ERROR] Invalid CanvasObject for updating cell: %T (expected *widget.Label)\n",
			obj)
		return
	}

	// switch id.Col {
	// case 0:
	// 	val = w.Location
	// case 1:
	// 	val = w.Period()[0].Format(common.TimestampFormatMinute)
	// case 2:
	// 	val = w.Period()[1].Format(common.TimestampFormatMinute)
	// case 3:
	// 	val = fmt.Sprintf("%s (%d)",
	// 		w.Event,
	// 		w.Level)
	// default:
	// 	g.log.Printf("[ERROR] Invalid column index: %d, valid range is 0 - 3\n",
	// 		id.Col)
	// 	return
	// }

	val = fmt.Sprintf("%s -> %s",
		w.Location,
		w.Event)

	lbl.SetText(val)
} // func (g *GUI) lstUpdate(id widget.TableCellID, obj fyne.CanvasObject)

// func (g *GUI) tblUpdate(id widget.TableCellID, obj fyne.CanvasObject) {
// 	g.log.Printf("[TRACE] tblUpdate -> Row = %d, Col = %d\n",
// 		id.Row,
// 		id.Col)

// 	var idx = id.Row

// 	if idx < 0 || idx >= len(g.warnings) {
// 		g.log.Printf("[ERROR] Invalid Row ID %d, valid range is 0 - %d\n",
// 			idx,
// 			len(g.warnings)-1)
// 		return
// 	}

// 	var (
// 		lbl *widget.Label
// 		val string
// 		w   = g.warnings[idx]
// 	)

// 	switch l := obj.(type) {
// 	case *widget.Label:
// 		lbl = l
// 	default:
// 		g.log.Printf("[ERROR] Invalid CanvasObject for updating cell: %T (expected *widget.Label)\n",
// 			obj)
// 		return
// 	}

// 	switch id.Col {
// 	case 0:
// 		val = w.Location
// 	case 1:
// 		val = w.Period()[0].Format(common.TimestampFormatMinute)
// 	case 2:
// 		val = w.Period()[1].Format(common.TimestampFormatMinute)
// 	case 3:
// 		val = fmt.Sprintf("%s (%d)",
// 			w.Event,
// 			w.Level)
// 	default:
// 		g.log.Printf("[ERROR] Invalid column index: %d, valid range is 0 - 3\n",
// 			id.Col)
// 		return
// 	}

// 	lbl.SetText(val)
// } // func (g *GUI) tblUpdate(id widget.TableCellID, obj fyne.CanvasObject)

func (g *GUI) refresh() {
	g.client.Refresh()
} // func (g *GUI) refresh()

func (g *GUI) wloop() {
	for w := range g.client.WarnQueue {
		g.log.Printf("[DEBUG] Received Warning from Client: %s -> %s\n",
			w.Location,
			w.Event)
		g.processWarning(w)
	}
} // func (g *GUI) wloop()

func (g *GUI) processWarning(w data.Warning) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.warnMap[w.ID] {
		g.log.Printf("[DEBUG] Append new Warning: %s -> %s\n",
			w.Location,
			w.Event)
		g.warnings = append(g.warnings, w)
		g.warnMap[w.ID] = true
	}
} // func (g *GUI) processWarning()
