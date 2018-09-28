# Reloader

[![GoDoc](https://godoc.org/gitlab.com/stop.start/go-chrome-remote-reload?status.svg)](https://godoc.org/gitlab.com/stop.start/go-chrome-remote-reload)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/stop.start/go-chrome-remote-reload)](https://goreportcard.com/report/gitlab.com/stop.start/go-chrome-remote-reload)

This library uses Chrome's remote debugging which allows to reload tabs remotely.

## Installation

Make sure you have a correctly configured Go environment. See [here](http://golang.org/doc/install.html) for instructions.  

Then, to install:  

```shell
go get gitlab.com/stop.start/go-chrome-remote-reload
```

## Getting started

RemoteConfig structure configures Chrome's remote debugging protocol and is used to open a new window as well as reload and close tabs.
The remote debugging protocol can be used also with Chromium or any Chrome-like browser supporting this protocol.
See [documetation](https://godoc.org/gitlab.com/stop.start/go-chrome-remote-reload#RemoteConfig) for details on the configuration.

The following will get the default config and open the browser on localhost:8080:

```go
package main

import(
    "gitlab.com/stop.start/go-chrome-remote-reload"
)

rc := RemoteConfigDefault()
rc.RemoteChrome()
```

To reload the only opened tab ReloadAllTabs is the easiest:

```go
rc.ReloadAllTabs()
```

## Examples

Use chromium instead of chrome:

```go
rc := RemoteConfigDefault()
rc.ExecName = "chromium"
rc.RemoteChrome()
```

Open two tabs and reload one of them:

```go
rc := reloader.NewRemoteConfig()

rc.Addr = "google.com"
rc.Port = 80
rc.RemoteChrome()

rc.Addr = "github.com"
rc.Port = 80 // not required since already set but more readable.
rc.RemoteChrome()

rc.ReloadTab("github.com/")
```

Reload tabs under same path:
Here only github.com/golang/go won't reload.

```go
rc := reloader.NewRemoteConfig()

rc.Addr = "github.com"
	rc.Port = 80
	rc.Route = "/golang/go"
	rc.RemoteChrome()

	rc.Addr = "github.com"
	rc.Port = 80
	rc.Route = "/golang/go/tree/master/misc"
	rc.RemoteChrome()

	rc.Addr = "github.com"
	rc.Port = 80
	rc.Route = "/golang/go/tree/master/test"
	rc.RemoteChrome()

rc.ReloadTab("tree/")
```

Close all tabs which also means closing Chrome :

```go
rc.CloseAllTabs()
```





