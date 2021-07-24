// /home/krylon/go/src/github.com/blicero/dwd/data/client.go
// -*- mode: go; coding: utf-8; -*-
// Created on 23. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-24 12:16:11 krylon>

// Package data implements the client to the DWD's web service, it fetches and
// processes the warning data.
package data

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/blicero/dwd/common"
	"github.com/blicero/dwd/logdomain"
)

const warnURL = "https://www.dwd.de/DWD/warnungen/warnapp/json/warnings.json"

// The response from the DWD's web service looks like this:
// warnWetter.loadWarnings({"time":1627052765000,"warnings":{},"vorabInformation":{},"copyright":"Copyright Deutscher Wetterdienst"});

var respPattern = regexp.MustCompile(`^warnWetter.loadWarnings\((.*)\);$`)

// Client implements the communication with the DWD's web service and the handling of the response.
type Client struct {
	log    *log.Logger
	client http.Client
}

// New creates a new Client. If proxy is a non-empty string, it is used as the
// URL of the proxy server to use for accessing the DWD's web service.
func New(proxy string) (*Client, error) {
	var (
		err error
		c   = new(Client)
	)

	if c.log, err = common.GetLogger(logdomain.Client); err != nil {
		return nil, err
	}

	c.client.Timeout = time.Second * 30

	if proxy != "" {
		var u *url.URL
		if u, err = url.Parse(proxy); err != nil {
			c.log.Printf("[ERROR] Cannot parse proxy URL %q: %s\n",
				proxy,
				err.Error())
			return nil, err
		}

		var pfunc = func(r *http.Request) (*url.URL, error) { return u, nil }

		switch t := c.client.Transport.(type) {
		case *http.Transport:
			t.Proxy = pfunc
		default:
			err = fmt.Errorf("Unexpected type for HTTP Client Transport: %T",
				c.client.Transport)
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return nil, err
		}

	}

	return c, nil
} // func New(proxy string) (*Client, error)

// FetchWarning fetches the warning data from the DWD's web service.
func (c *Client) FetchWarning() ([]byte, error) {
	var (
		err  error
		n    int64
		res  *http.Response
		buf  bytes.Buffer
		body []byte
	)

	if res, err = c.client.Get(warnURL); err != nil {
		c.log.Printf("[ERROR] Failed to fetch %q: %s\n",
			warnURL,
			err.Error())
	}

	c.log.Printf("[DEBUG] Response for %q: %s\n",
		warnURL,
		res.Status)

	defer res.Body.Close() // nolint: errcheck

	if n, err = io.Copy(&buf, res.Body); err != nil {
		c.log.Printf("[ERROR] Cannot read response Body for %q: %s\n",
			warnURL,
			err.Error())
		return nil, err
	}

	body = buf.Bytes()

	c.log.Printf("[DEBUG] Response from %s: %s (%d bytes of pure %s)\n",
		warnURL,
		res.Status,
		n,
		res.Header.Get("Content-Type"))

	var match [][]byte

	if match = respPattern.FindSubmatch(body[:n]); match == nil {
		err = fmt.Errorf("Cannot parse response from %q: %q",
			warnURL,
			body[:n])
		c.log.Printf("[ERROR] %s\n", err.Error())
		return nil, err
	}

	var data = match[1]

	// c.log.Printf("[DEBUG] Received response from DWD: %s\n",
	// 	data)

	return data, nil
} // func (c *Client) FetchWarning() ([]byte, error)
