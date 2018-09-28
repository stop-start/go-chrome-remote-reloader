package reloader_test

import (
	"testing"

	"time"

	"gitlab.com/stop.start/go-chrome-remote-reload"
)

/*
	Because everything happens in Chrome browser, regular testing is not really
	an option.
	All the tests here rely on visual: they all open a chrome window with various
	scenarios like reload one tab only, closing all tabs, etc.
*/

func TestRemoteChromeOpenAndCloseBrowser(t *testing.T) {
	/*
		Should open the browser with one tab without nothing in it (unsless
		there's an active server on port 8080)
	*/
	rc := reloader.NewRemoteConfig()
	rc.ExecName = "chromium"
	rc.RemoteChrome()
	time.Sleep(time.Second * 2)
	rc.CloseAllTabs()
	time.Sleep(time.Second * 2)
}

func TestRemoteChromeReloadOnlyTab(t *testing.T) {
	/*
		Should open the browser with website http://randomword.com
		After a few second the tab will reload and should show a different word
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
	time.Sleep(time.Second * 2)
	rc.CloseAllTabs()
	time.Sleep(time.Second * 2)
}

func TestRemoteChromeReloadAllTabs(t *testing.T) {
	/*
		Should open the browser with websites http://randomword.com and
		http://randomfactgenerator.net.
		After a few seconds the 2 tabs will reload and should show a different
		content.
	*/
	rc := reloader.NewRemoteConfig()
	rc.ExecName = "chromium"
	rc.Addr = "randomword.com"
	rc.Port = 80
	rc.RemoteChrome()

	rc.Addr = "randomfactgenerator.net"
	rc.Port = 80
	rc.RemoteChrome()

	time.Sleep(time.Second * 4)
	err := rc.ReloadAllTabs()
	if err != nil {
		t.Errorf("Error while reloading: %s", err)
	}
	time.Sleep(time.Second * 4)
	rc.CloseAllTabs()
	time.Sleep(time.Second * 2)
}

func TestRemoteChromeReloadGroupTabs(t *testing.T) {
	/*
		Should open the browser with websites http://randomword.com and
		some page in the golang repository in github.com.
		After a few second the tabs under /golang/go/tree will reload.
	*/
	rc := reloader.NewRemoteConfig()
	rc.ExecName = "chromium"
	rc.Addr = "randomword.com"
	rc.Port = 80
	rc.RemoteChrome()

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

	time.Sleep(time.Second * 6)
	err := rc.ReloadTabGroup("tree/")
	if err != nil {
		t.Errorf("Error while reloading: %s", err)
	}
	time.Sleep(time.Second * 6)
	rc.CloseAllTabs()
	time.Sleep(time.Second * 3)
}

func TestRemoteChromeCloseOneTab(t *testing.T) {
	/*
		Should open the browser with websites http://randomword.com and
		http://randomfactgenerator.net.
		After a few seconds the http://randomfactgenerator.net tab
		should close.
	*/
	rc := reloader.NewRemoteConfig()
	rc.ExecName = "chromium"
	rc.Addr = "randomword.com"
	rc.Port = 80
	rc.RemoteChrome()

	rc.Addr = "randomfactgenerator.net"
	rc.Port = 80
	rc.RemoteChrome()

	time.Sleep(time.Second * 4)
	err := rc.CloseTab("randomfactgenerator.net/")
	if err != nil {
		t.Errorf("Error while reloading: %s", err)
	}
	time.Sleep(time.Second * 4)
	rc.CloseAllTabs()
	time.Sleep(time.Second * 2)
}

func TestRemoteChromeCloseGroupTabs(t *testing.T) {
	/*
		Should open the browser with websites http://randomword.com and
		some page in the golang repository in github.com.
		After a few second the tabs under /golang/go/tree should close.
	*/
	rc := reloader.NewRemoteConfig()
	rc.ExecName = "chromium"
	rc.Addr = "randomword.com"
	rc.Port = 80
	rc.RemoteChrome()

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

	time.Sleep(time.Second * 4)
	err := rc.CloseTabGroup("tree/")
	if err != nil {
		t.Errorf("Error while reloading: %s", err)
	}
	time.Sleep(time.Second * 4)
	rc.CloseAllTabs()
	time.Sleep(time.Second * 2)
}
