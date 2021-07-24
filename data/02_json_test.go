// /home/krylon/go/src/github.com/blicero/dwd/data/02_json_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 24. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-24 12:14:49 krylon>

package data

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestParseJSON(t *testing.T) {
	const path = "testdata/warnings.1.json"
	var (
		err  error
		fh   *os.File
		buf  bytes.Buffer
		info WeatherInfo
	)

	if fh, err = os.Open(path); err != nil {
		t.Fatalf("Error opening %s: %s",
			path,
			err.Error())
	}

	defer fh.Close() // nolint: errcheck

	if _, err = io.Copy(&buf, fh); err != nil {
		t.Fatalf("Error reading %s: %s",
			path,
			err.Error())
	} else if err = json.Unmarshal(buf.Bytes(), &info); err != nil {
		t.Fatalf("Error parsing %s: %s",
			path,
			err.Error())
	}

} // func TestParseJSON(t *testing.T)
