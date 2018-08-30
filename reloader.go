package reloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"golang.org/x/net/websocket"
)

// Parameters  struct ...
type Parameters struct {
	IgnoreCache bool `json:"ignoreCache"`
}

// RefreshJSON struct ...
type RefreshJSON struct {
	ID     uint16     `json:"id"`
	Method string     `json:"method"`
	Params Parameters `json:"params"`
}

// ChromeTab is the json returned from chrome remote debugging api
type ChromeTab struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	Description          string `json:"description"`
}

// RemoteConfig is the configuration used to start chrome remote debugging
type RemoteConfig struct {
	// ExecName is the name of the chrome executable (e.g "google-chrome", "chromium")
	ExecName string
	// Port is the chrome remote debugging port, usually 9222
	Port int
	// UserDataDir is the folder for chrome user profile
	UserDataDir string
	// OriginAddr is the address to be reloaded
	OriginAddr string
	// OriginPort is the port for the address to be reloaded
	OriginPort int
}

// RemoteChrome starts a new chrome remote debugging session
func RemoteChrome(rc *RemoteConfig) chan error {
	return rc.remoteChrome()
}

// RemoteChromeDefault starts a new chrome remote debugging session with default port and user data directory
func RemoteChromeDefault() error {
	rc := &RemoteConfig{
		ExecName:    "chromium", //TODO: change to google-chrome
		Port:        9222,
		UserDataDir: "~/.chrome-remote-profile",
		OriginAddr:  "localhost",
		OriginPort:  8080,
	}
	return rc.remoteChrome()
}

func (rc *RemoteConfig) remoteChrome() chan error {
	cmd := exec.Command(
		rc.ExecName,
		"--remote-debugging-port", rc.Port,
		"--user-data-dir", rc.UserDataDir,
		fmt.Sprintf("http://%s:%d", rc.OriginAddr, rc.OriginPort),
	)

	errChan := make(chan error)

	go func() {
		if err := cmd.Run(); err != nil {
			errChan <- err
		}
		errChan <- nil
	}()

	return errChan
}

// ReloadAllTabs will reload all opened tabs
func (rc *RemoteConfig) ReloadAllTabs() error {
	tabs, err := getTabs(rc.OriginAddr, rc.OriginPort)
	if err != nil {
		return err
	}
	var err error
	for _, tab := range tabs {
		e := reloadTab(tab)
		if e != nil {
			err = e
		}
	}
	return e
}

// ReloadTab reloads one chrome tab by checking if the tab URL has route as suffix
func ReloadTab(route string) error {
	tabs, err := getTabs(rc.OriginAddr, rc.OriginPort)
	if err != nil {
		return err
	}

	for _, tab := range tabs {
		if strings.HasSuffix(tab.URL, route) {
			err := reloadTab(tab)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

// ReloadTabGroup reloads a group of chrome tabs by checking if the tab URL contains the subroute
func ReloadTabGroup(subroute string) error {
	tabs, err := getTabs(rc.OriginAddr, rc.OriginPort)
	if err != nil {
		return err
	}

	var err error
	for _, tab := range tabs {
		if strings.Contains(tab.URL, subroute) {
			e := reloadTab(tab)
			if e != nil {
				err = e
			}
		}
	}
	return err

}

func reloadTab(tab ChromeTab) error {
	url := tab.WebSocketDebuggerURL
	origin := tab.URL
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return err
	}

	jsonStruct := &RefreshJSON{ID: 0, Method: "Page.reload", Params: Parameters{IgnoreCache: true}}
	jsonString, _ := json.Marshal(jsonStruct)
	_, err = ws.Write(jsonString)
	if err != nil {
		return err
	}
	return nil
}

func getTabs(addr string, port int) ([]ChromeTab, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/json", addr, port))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tabs []ChromeTab
	buffer, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buffer, &tabs)
	if err != nil {
		return nil, err
	}
	return tabs, nil
}
