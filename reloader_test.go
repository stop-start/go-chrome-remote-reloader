package reloader_test

import (
	"testing"

	"time"

	"gitlab.com/stop.start/go-chrome-remote-reload"
)

func TestRemoteChromeOpenBrowser(t *testing.T) {
	/*
		Should open the browser with one tab without nothing in it (unsless
		there's an active server on port 8080)
	*/
	rc := reloader.NewRemoteConfig()
	rc.ExecName = "chromium"
	rc.RemoteChrome()
}

func TestRemoteChromeReloadOnlyTab(t *testing.T) {
	/*
		Should open the browser with website http://randomword.com
		After 2 second the tab will reload and should show a different word
	*/
	rc := reloader.NewRemoteConfig()
	rc.ExecName = "chromium"
	rc.Addr = "randomword.com"
	rc.Port = 80
	rc.RemoteChrome()
	time.Sleep(time.Second * 2)
	err := rc.ReloadTab(rc.Route)
	if err != nil {
		t.Errorf("Error while reloading: %s", err)
	}
	rc
}
