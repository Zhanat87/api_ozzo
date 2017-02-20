// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

// This file implements an http.Client with request timeouts set by command
// line flags.

package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/memcache"

	"github.com/golang/gddo/httputil"
)

var (
	dialTimeout    = flag.Duration("dial_timeout", 5*time.Second, "Timeout for dialing an HTTP connection.")
	requestTimeout = flag.Duration("request_timeout", 20*time.Second, "Time out for roundtripping an HTTP request.")
)

func newHTTPClient() *http.Client {
	t := newCacheTransport()
	t.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   *dialTimeout,
			KeepAlive: *requestTimeout / 2,
		}).Dial,
		ResponseHeaderTimeout: *requestTimeout / 2,
		TLSHandshakeTimeout:   *requestTimeout / 2,
	}
	return &http.Client{
		// Wrap the cached transport with GitHub authentication.
		Transport: httputil.NewAuthTransport(t),
		Timeout:   *requestTimeout,
	}
}

func newCacheTransport() *httpcache.Transport {
	// host and port are set by GAE Flex runtime, can be left blank locally.
	host := os.Getenv("MEMCACHE_PORT_11211_TCP_ADDR")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("MEMCACHE_PORT_11211_TCP_PORT")
	if port == "" {
		port = "11211"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	return httpcache.NewTransport(
		memcache.New(addr),
	)
}
