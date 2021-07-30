// /home/krylon/go/src/github.com/blicero/dwd/main.go
// -*- mode: go; coding: utf-8; -*-
// Created on 24. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-29 22:37:16 krylon>

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blicero/dwd/common"
	"github.com/blicero/dwd/data"
	"github.com/blicero/dwd/ui/fyne"
	"github.com/blicero/dwd/ui/gtk3"
)

var locations = []string{
	"(?i)bielefeld",
	"(?i)gütersloh",
	"(?i)paderborn",
	"(?i)osnabrück",
}

const interval = time.Second * 60

func main() {
	var (
		err    error
		ui     string
		client *data.Client
		sigQ   = make(chan os.Signal)
	)

	flag.StringVar(&ui, "interface", "gtk", "The user interface to present (currently: cli, fyne, gtk)")
	flag.Parse()

	if len(flag.Args()) > 0 {
		locations = append(locations, flag.Args()...)
	}

	if client, err = data.New("", locations...); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Client: %s\n",
			err.Error())
		os.Exit(1)
	}

	client.Start()

	signal.Notify(sigQ, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	switch ui {
	case "fyne":
		fmt.Println("Starting fyne GUI")

		var g *fyne.GUI

		if g, err = fyne.MakeGUI(client); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create GUI: %s\n",
				err.Error())
			os.Exit(1)
		}

		g.ShowAndRun()
	case "gtk":
		fmt.Println("Starting gtk GUI")

		var g *gtk3.GUI

		if g, err = gtk3.MakeGUI(client); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create GUI: %s\n",
				err.Error())
			os.Exit(1)
		}

		g.ShowAndRun()
	case "cli":
		fmt.Println("Starting command line interface")

		var (
			ticker = time.NewTicker(time.Second * 10)
			idMap  = make(map[int64]bool)
		)

		for client.IsActive() {
			select {
			case <-ticker.C:
				// ...
			case <-sigQ:
				fmt.Println("Quitting because of signal")
				//os.Exit(0)
				client.Stop()
			case w := <-client.WarnQueue:
				if idMap[w.ID] {
					continue
				}
				idMap[w.ID] = true
				var p = w.Period()
				fmt.Printf("Warnung %d für %s von %s bis %s: %s\n",
					w.ID,
					w.Location,
					p[0].Format(common.TimestampFormatMinute),
					p[1].Format(common.TimestampFormatMinute),
					w.Event)
			}
		}
	}
}
