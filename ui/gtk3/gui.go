// /home/krylon/go/src/github.com/blicero/dwd/ui/gtk3/gui.go
// -*- mode: go; coding: utf-8; -*-
// Created on 29. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-08-01 14:29:41 krylon>

// Package gtk3 provides a gtk3-based GUI.
// Getting gotk3 to compile cleanly across multiple platforms has been kind of
// a pain so far, but if I can get it work, it should be rather nice.
package gtk3

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/blicero/dwd/common"
	"github.com/blicero/dwd/data"
	"github.com/blicero/dwd/logdomain"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const refInterval = time.Second * 10

const (
	cidLocation = iota
	cidEvent
	cidLevel
	cidStart
	cidEnd
)

// GUI wraps the UI with all its assorted state.
type GUI struct {
	log         *log.Logger
	lock        sync.RWMutex
	warnMap     map[string]data.Warning
	updateStamp time.Time
	renderStamp time.Time
	client      *data.Client
	win         *gtk.Window
	box         *gtk.Box
	refB        *gtk.Button
	wView       *gtk.TreeView
	wStore      *gtk.ListStore
	wSort       *gtk.TreeModelSort
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
	} else if g.wStore, err = gtk.ListStoreNew(
		glib.TYPE_STRING, // Location
		glib.TYPE_STRING, // Event
		glib.TYPE_INT,    // Level
		glib.TYPE_STRING, // Start
		glib.TYPE_STRING, // End
	); err != nil {
		g.log.Printf("[ERROR] Cannot create ListStore: %s\n",
			err.Error())
		return nil, err
	} else if g.wSort, err = gtk.TreeModelSortNew(g.wStore); err != nil {
		g.log.Printf("[ERROR] Cannot create TreeModelSort: %s\n",
			err.Error())
		return nil, err
	} else if g.wView, err = gtk.TreeViewNewWithModel(g.wSort); err != nil {
		g.log.Printf("[ERROR] Cannot create TreeView: %s\n",
			err.Error())
		return nil, err
	} else if g.wScr, err = gtk.ScrolledWindowNew(nil, nil); err != nil {
		g.log.Printf("[ERROR] Cannot create ScrolledWindow: %s\n",
			err.Error())
		return nil, err
	}

	g.win.Stick()

	g.wSort.SetDefaultSortFunc(g.cmpWarnings)

	var cLocation, cEvent, cLevel, cStart, cEnd *gtk.TreeViewColumn

	if cLocation, err = createCol("Ort", cidLocation); err != nil {
		g.log.Printf("[ERROR] Cannot create TreeViewColumn for Location: %s\n",
			err.Error())
		return nil, err
	} else if cEvent, err = createCol("Ereignis", cidEvent); err != nil {
		g.log.Printf("[ERROR] Cannot create TreeViewColumn for Event: %s\n",
			err.Error())
		return nil, err
	} else if cLevel, err = createCol("Warnstufe", cidLevel); err != nil {
		g.log.Printf("[ERROR] Cannot create TreeViewColumn for Warnstufe: %s\n",
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
	cLevel.SetVisible(true)
	cStart.SetVisible(true)
	cEnd.SetVisible(true)

	g.wView.AppendColumn(cLocation)
	g.wView.AppendColumn(cEvent)
	g.wView.AppendColumn(cLevel)
	g.wView.AppendColumn(cStart)
	g.wView.AppendColumn(cEnd)

	g.warnMap = make(map[string]data.Warning)

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

	if g.renderStamp.Before(g.updateStamp) {
		g.log.Printf("[DEBUG] Update Label (%d warnings)\n", len(g.warnMap))
		g.wStore.Clear()

		for _, w := range g.warnMap {
			var (
				err error
				p   = w.Period()
			)

			var iter = g.wStore.Append()

			if err = g.wStore.Set(iter,
				[]int{0, 1, 2, 3, 4},
				[]interface{}{
					w.Location,
					w.Event,
					w.Level,
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

	if _, ok := g.warnMap[w.GetUniqueID()]; !ok {
		g.log.Printf("[DEBUG] Append new Warning: %s -> %s\n",
			w.Location,
			w.Event)
		g.warnMap[w.GetUniqueID()] = w
		g.updateStamp = time.Now()
	}
} // func (g *GUI) processWarnings(w data.Warning)

func (g *GUI) cmpWarnings(m *gtk.TreeModel, a, b *gtk.TreeIter) int {
	var (
		err    error
		v1, v2 *glib.Value
		l1, l2 string
	)

	if v1, err = m.GetValue(a, cidLocation); err != nil {
		g.log.Printf("[ERROR] Cannot get TreeModel value %p: %s\n",
			a,
			err.Error())
		return 0
	} else if v2, err = m.GetValue(b, cidLocation); err != nil {
		g.log.Printf("[ERROR] Cannot get TreeModel value %p: %s\n",
			b,
			err.Error())
		return 0
	} else if l1, err = v1.GetString(); err != nil {
		g.log.Printf("[ERROR] Cannot get string value of %p: %s\n",
			a,
			err.Error())
		return 0
	} else if l2, err = v2.GetString(); err != nil {
		g.log.Printf("[ERROR] Cannot get string value of %p: %s\n",
			b,
			err.Error())
		return 0
	}

	switch strings.Compare(l1, l2) {
	case -1:
		return -1
	case 0:
		if v1, err = m.GetValue(a, 2); err != nil {
			g.log.Printf("[ERROR] Cannot get string values of %p/2: %s\n",
				a,
				err.Error())
			return 0
		} else if v2, err = m.GetValue(b, 2); err != nil {
			g.log.Printf("[ERROR] Cannot get string values of %p/2: %s\n",
				b,
				err.Error())
			return 0
		} else if l1, err = v1.GetString(); err != nil {
			g.log.Printf("[ERROR] Cannot get String value of %p: %s\n",
				a,
				err.Error())
			return 0
		} else if l2, err = v2.GetString(); err != nil {
			g.log.Printf("[ERROR] Cannot get String value of %p: %s\n",
				b,
				err.Error())
			return 0
		}

		return strings.Compare(l1, l2)
	case 1:
		return 1
	default:
		return 0
	}
} // func cmpWarnings(m *gtk.TreeModel, a, b, *gtk.TreeIter) int
