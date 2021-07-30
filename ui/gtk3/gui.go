// /home/krylon/go/src/github.com/blicero/dwd/ui/gtk3/gui.go
// -*- mode: go; coding: utf-8; -*-
// Created on 29. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-30 18:01:18 krylon>

// Package gtk3 provides a gtk3-based GUI.
// Getting gotk3 to compile cleanly across multiple platforms has been kind of
// a pain so far, but if I can get it work, it should be rather nice.
package gtk3

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/blicero/dwd/common"
	"github.com/blicero/dwd/data"
	"github.com/blicero/dwd/logdomain"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const refInterval = time.Second * 5

const (
	cidLocation = iota
	cidEvent
	cidStart
	cidEnd
)

// GUI wraps the UI with all its assorted state.
type GUI struct {
	log         *log.Logger
	lock        sync.RWMutex
	warnMap     map[int64]bool
	warnings    []data.Warning
	updateStamp time.Time
	renderStamp time.Time
	client      *data.Client
	win         *gtk.Window
	box         *gtk.Box
	refB        *gtk.Button
	wView       *gtk.TreeView
	wStore      *gtk.ListStore
	wScr        *gtk.ScrolledWindow
}

// MakeGUI creates a new Gtk3 GUI.
func MakeGUI(c *data.Client) (*GUI, error) {
	var (
		err error
		g   = &GUI{client: c}
	)

	if g.log, err = common.GetLogger(logdomain.GUI); err != nil {
		return nil, err
	}

	gtk.Init(nil)

	if g.win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL); err != nil {
		g.log.Printf("[ERROR] Cannot create Gtk window: %s\n",
			err.Error())
		return nil, err
	} else if g.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1); err != nil {
		g.log.Printf("[ERROR] Cannot create Gtk Box: %s\n",
			err.Error())
		return nil, err
	} else if g.refB, err = gtk.ButtonNewWithMnemonic("_Aktualisieren"); err != nil {
		g.log.Printf("[ERROR] Cannot create Gtk Button: %s\n",
			err.Error())
		return nil, err
	} else if g.wStore, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING); err != nil {
		g.log.Printf("[ERROR] Cannot create ListStore: %s\n",
			err.Error())
		return nil, err
	} else if g.wView, err = gtk.TreeViewNewWithModel(g.wStore); err != nil {
		g.log.Printf("[ERROR] Cannot create TreeView: %s\n",
			err.Error())
		return nil, err
	} else if g.wScr, err = gtk.ScrolledWindowNew(nil, nil); err != nil {
		g.log.Printf("[ERROR] Cannot create ScrolledWindow: %s\n",
			err.Error())
		return nil, err
	}

	var cLocation, cEvent, cStart, cEnd *gtk.TreeViewColumn

	if cLocation, err = createCol("Ort", cidLocation); err != nil {
		g.log.Printf("[ERROR] Cannot create TreeViewColumn for Location: %s\n",
			err.Error())
		return nil, err
	} else if cEvent, err = createCol("Ereignis", cidEvent); err != nil {
		g.log.Printf("[ERROR] Cannot create TreeViewColumn for Event: %s\n",
			err.Error())
		return nil, err
	} else if cStart, err = createCol("Von", cidStart); err != nil {
		g.log.Printf("[ERROR] Cannot create TreeViewColumn for Start time: %s\n",
			err.Error())
		return nil, err
	} else if cEnd, err = createCol("Bis", cidEnd); err != nil {
		g.log.Printf("[ERROR] Cannot create TreeViewColumn for End time: %s\n",
			err.Error())
		return nil, err
	}

	cLocation.SetVisible(true)
	cEvent.SetVisible(true)
	cStart.SetVisible(true)
	cEnd.SetVisible(true)

	g.wView.AppendColumn(cLocation)
	g.wView.AppendColumn(cEvent)
	g.wView.AppendColumn(cStart)
	g.wView.AppendColumn(cEnd)

	g.warnMap = make(map[int64]bool)

	g.win.Connect("destroy", gtk.MainQuit)
	g.refB.Connect("clicked", g.refresh)

	g.wScr.Add(g.wView)
	g.box.PackStart(g.refB, false, false, 0)
	g.box.PackStart(g.wScr, true, true, 0)
	g.win.Add(g.box)
	g.wScr.SetSizeRequest(480, 360)
	g.win.SetSizeRequest(480, 360)

	return g, nil
} // func MakeGUI(c *data.Client) (*GUI, error)

func (g *GUI) ShowAndRun() {
	g.win.ShowAll()
	glib.TimeoutAdd(uint(refInterval.Milliseconds()), g.timeoutHandler)
	go g.wloop()
	gtk.Main()
} // func (g *GUI) ShowAndRun()

func (g *GUI) timeoutHandler() bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	// g.log.Println("[DEBUG] timeoutHandler")

	if g.renderStamp.Before(g.updateStamp) {
		g.log.Printf("[DEBUG] Update Label (%d warnings)\n", len(g.warnings))
		g.wStore.Clear()

		for _, w := range g.warnings {
			var (
				err error
				p   = w.Period()
			)

			var iter = g.wStore.Append()

			if err = g.wStore.Set(iter,
				[]int{0, 1, 2, 3},
				[]interface{}{
					w.Location,
					w.Event,
					p[0].Format(common.TimestampFormatMinute),
					p[1].Format(common.TimestampFormatMinute),
				}); err != nil {
				g.log.Printf("[ERROR] Cannot add Warning to TreeView: %s\n",
					err.Error())
			}
		}

		g.renderStamp = time.Now()
	}

	return true
} // func (g *GUI) timeoutHandler() bool

func (g *GUI) refresh() {
	g.log.Println("[DEBUG] Refresh THIS!")
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

		// g.wLbl.SetText(fmt.Sprintf("%s -> %s",
		// 	w.Location,
		// 	w.Event))

		sort.Sort(data.WarningList(g.warnings))
		g.updateStamp = time.Now()
	}
} // func (g *GUI) processWarnings(w data.Warning)
