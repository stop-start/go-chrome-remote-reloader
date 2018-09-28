# Reloader

[![GoDoc](https://godoc.org/gitlab.com/stop.start/go-chrome-remote-reload?status.svg)](https://godoc.org/gitlab.com/stop.start/go-chrome-remote-reload)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/stop.start/go-chrome-remote-reload)](https://goreportcard.com/report/gitlab.com/stop.start/go-chrome-remote-reload)

This library uses Chrome's remote debugging which allows to reload tabs remotely.

## Installation

Make sure you have a correctly configured Go environment. See [here](http://golang.org/doc/install.html) for instructions.  

Then to install:  

```shell
go get gitlab.com/stop.start/go-chrome-remote-reload
```

## Getting started

### Default session

RemoteChromeDefault method allows to start a new Chrome session with default configuration.  
The code below will open the Chrome browser on localhost:8080. 

```go
package main

import(
    "gitlab.com/stop.start/go-chrome-remote-reload"
)

rc, _, err := RemoteChromeDefault() 
```

The following code will reload all opened tabs (here just the one created):

```go
rc.ReloadAllTabs()
```

### Custom configuration

RemoteConfig structure configures Chrome's remote debugging protocol and is used to start a new session as well as reload tabs.

The remote debugging protocol can be used also with Chromium or any Chrome-like browser supporting this protocol.

The following will get the default config and change the browser executable:

```go
package main

import(
    "gitlab.com/stop.start/go-chrome-remote-reload"
)

rc := RemoteConfigDefault() 
rc.ExecName = "chromium"
```




