// /home/krylon/go/src/github.com/blicero/dwd/main.go
// -*- mode: go; coding: utf-8; -*-
// Created on 24. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-26 12:42:25 krylon>

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/blicero/dwd/common"
	"github.com/blicero/dwd/data"
)

var locations = []string{
	"(?i)bielefeld",
	"(?i)gütersloh",
	"(?i)paderborn",
	"(?i)osnabrück",
}

const interval = time.Second * 60

func main() {
	fmt.Printf("%s: Implement me!\n", common.AppName)

	var (
		err      error
		client   *data.Client
		warnings []data.Warning
	)

	flag.Parse()

	if len(flag.Args()) > 0 {
		locations = append(locations, flag.Args()...)
	}

	if client, err = data.New("", locations...); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Client: %s\n",
			err.Error())
		os.Exit(1)
	}

	for {
		if warnings, err = client.GetWarnings(); err != nil {
			fmt.Fprintf(os.Stderr, "Error getting warnings: %s\n",
				err.Error())
			goto WAIT
		}

		for _, w := range warnings {
			fmt.Printf("Warning for %s: %s\n%s\n\n%s\n====================================\n",
				w.Location,
				w.Event,
				w.Description,
				w.Instruction)
		}

	WAIT:
		time.Sleep(interval)
	}
}
