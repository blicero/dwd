// /home/krylon/go/src/github.com/blicero/dwd/ui/fyne/stub.go
// -*- mode: go; coding: utf-8; -*-
// Created on 29. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-29 21:41:07 krylon>

// +build 386

package fyne

import (
	"errors"
	"fmt"
	"os"

	"github.com/blicero/dwd/data"
)

type GUI int

func MakeGUI(c *data.Client) (*GUI, error) {
	return nil, errors.New("Fyne is not available on 386")
} // func MakeGUI(c *data.Client) (*GUI, error)

func (g *GUI) ShowAndRun() {
	fmt.Fprintf(os.Stderr, "Fyne is not available on 386\n")
} // func (g *GUI) ShowAndRun()
