// /home/krylon/go/src/github.com/blicero/dwd/data/01_client_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 23. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-23 21:46:15 krylon>

package data

import "testing"

var c *Client

func TestClientCreate(t *testing.T) {
	var err error

	if c, err = New(""); err != nil {
		c = nil
		t.Fatalf("Failed to create Client: %s",
			err.Error())
	}
} // func TestClientCreate(t *testing.T)

func TestClientFetch(t *testing.T) {
	if c == nil {
		t.SkipNow()
	}

	var (
		data []byte
		err  error
	)

	if data, err = c.FetchWarning(); err != nil {
		t.Errorf("Failed to fetch weather warnings data: %s", err.Error())
	} else {
		t.Logf("Fetched data: %s", data)
	}
} // func TestClientFetch(t *testing.T)
