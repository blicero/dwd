// /home/krylon/go/src/github.com/blicero/dwd/data/client.go
// -*- mode: go; coding: utf-8; -*-
// Created on 23. 07. 2021 by Benjamin Walkenhorst
// (c) 2021 Benjamin Walkenhorst
// Time-stamp: <2021-07-26 19:10:58 krylon>

// Package data implements the client to the DWD's web service, it fetches and
// processes the warning data.
package data

import (
	"bytes"
	"encoding/json"
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

var respPattern = regexp.MustCompile(`^warnWetter.loadWarnings\((.*)\);`)

// Client implements the communication with the DWD's web service and the handling of the response.
type Client struct {
	log       *log.Logger
	client    http.Client
	locations []*regexp.Regexp
}

// New creates a new Client. If proxy is a non-empty string, it is used as the
// URL of the proxy server to use for accessing the DWD's web service.
func New(proxy string, locations ...string) (*Client, error) {
	var (
		err error
		c   = new(Client)
	)

	if c.log, err = common.GetLogger(logdomain.Client); err != nil {
		return nil, err
	}

	c.client.Timeout = time.Second * 90

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

	c.locations = make([]*regexp.Regexp, 0, len(locations))

	for _, l := range locations {
		var r *regexp.Regexp

		c.log.Printf("[DEBUG] Add regexp %s\n", l)

		if r, err = regexp.Compile(l); err != nil {
			c.log.Printf("[ERROR] Cannot compile Regexp %q: %s\n",
				l,
				err.Error())
			return nil, err
		}

		c.locations = append(c.locations, r)
	}

	c.log.Printf("[DEBUG] Client has %d regular expressions for matching locations\n",
		len(c.locations))

	return c, nil
} // func New(proxy string) (*Client, error)

// ProcessWarnings parses the warnings returned by the DWD's web service and
// returns a list of all the warnings that are relevant to us.
func (c *Client) ProcessWarnings(raw []byte) ([]Warning, error) {
	var (
		err  error
		info WeatherInfo
	)

	if err = json.Unmarshal(raw, &info); err != nil {
		c.log.Printf("[ERROR] Cannot parse JSON data: %s\n%s\n",
			err.Error(),
			raw)
		return nil, err
	}

	var list = make([]Warning, 0, len(c.locations)*2)

	for _, i := range info.Warnings {
	W_ITEM:
		for _, w := range i {
			for _, l := range c.locations {
				if m := l.FindString(w.Location); m != "" {
					// c.log.Printf("[DEBUG] Found Match for %s: %s\n",
					// 	l,
					// 	w.Location)
					list = append(list, w)
					continue W_ITEM
				}
			}
		}
	}

	for _, i := range info.PrelimWarnings {
	V_ITEM:
		for _, w := range i {
			for _, l := range c.locations {
				if m := l.FindString(w.Location); m != "" {
					c.log.Printf("[DEBUG] Found Match for %s: %s\n",
						l,
						w.Location)
					list = append(list, w)
					continue V_ITEM
				}
			}
		}
	}

	return list, nil
} // func (c *Client) ProcessWarnings(raw []byte) ([]Warning, error)

// FetchWarning fetches the warning data from the DWD's web service.
func (c *Client) FetchWarning() ([]byte, error) {
	var (
		err error
		res *http.Response
		buf bytes.Buffer
	)

	if res, err = c.client.Get(warnURL); err != nil {
		c.log.Printf("[ERROR] Failed to fetch %q: %s\n",
			warnURL,
			err.Error())
	}

	defer res.Body.Close() // nolint: errcheck

	if res.StatusCode != 200 {
		c.log.Printf("[DEBUG] Response for %q: %s\n",
			warnURL,
			res.Status)
		return nil, fmt.Errorf("HTTP Request to %q failed: %s",
			warnURL,
			res.Status)
	} else if _, err = io.Copy(&buf, res.Body); err != nil {
		c.log.Printf("[ERROR] Cannot read response Body for %q: %s\n",
			warnURL,
			err.Error())
		return nil, err
	}

	var body = buf.Bytes()

	// c.log.Printf("[DEBUG] Response from %s: %s (%d bytes of pure %s)\n",
	// 	warnURL,
	// 	res.Status,
	// 	n,
	// 	res.Header.Get("Content-Type"))

	var match [][]byte

	if match = respPattern.FindSubmatch(body); match == nil {
		err = fmt.Errorf("Cannot parse response from %q: %q",
			warnURL,
			body)
		c.log.Printf("[ERROR] %s\n", err.Error())
		return nil, err
	}

	var data = match[1]

	// c.log.Printf("[DEBUG] Received response from DWD: %s\n",
	// 	data)

	return data, nil
} // func (c *Client) FetchWarning() ([]byte, error)

// GetWarnings loads the current warnings from the DWD and returns all warnings
// matching its list of locations.
func (c *Client) GetWarnings() ([]Warning, error) {
	var (
		err      error
		rawData  []byte
		warnings []Warning
	)

	if rawData, err = c.FetchWarning(); err != nil {
		c.log.Printf("[ERROR] Failed to fetch data from DWD: %s\n",
			err.Error())
		return nil, err
	} else if warnings, err = c.ProcessWarnings(rawData); err != nil {
		c.log.Printf("[ERROR] Failed to process Warnings: %s\n",
			err.Error())
		return nil, err
	}

	return warnings, nil
} // func (c *Client) GetWarnings() ([]Warning, error)
