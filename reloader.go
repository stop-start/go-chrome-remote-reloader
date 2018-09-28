package reloader

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"golang.org/x/net/websocket"
)

// RemoteConfig is the configuration used to start chrome remote debugging
type RemoteConfig struct {
	// ExecName is the name of the chrome executable (e.g "google-chrome", "chromium")
	ExecName string
	// Port is the chrome remote debugging port, usually 9222
	Port int
	// UserDataDir is the folder for chrome user profile
	UserDataDir string
	// OriginAddr is the base address to be opened
	OriginAddr string
	// OriginPort is the port for the base address to be opened
	OriginPort int
	//OriginRoute path to be opened
	OriginRoute string
}

// json returned from chrome remote debugging api
type chromeTab struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	Description          string `json:"description"`
}

// sent to chrome debugging port to reload tabs
type refreshJSON struct {
	ID     uint16     `json:"id"`
	Method string     `json:"method"`
	Params parameters `json:"params"`
}

type parameters struct {
	IgnoreCache bool `json:"ignoreCache"`
}

// RemoteConfigDefault returns a default config for RemoteChrome
func RemoteConfigDefault() *RemoteConfig {
	return &RemoteConfig{
		ExecName:    "google-chrome",
		Port:        9222,
		UserDataDir: "/tmp/.chrome-remote-profile",
		OriginAddr:  "localhost",
		OriginPort:  8080,
		OriginRoute: "",
	}
}

// RemoteChrome starts a new chrome remote debugging session
func RemoteChrome(rc *RemoteConfig) (context.CancelFunc, error) {
	return rc.remoteChrome()
}

// RemoteChromeDefault starts a new chrome remote debugging session with default port and user data directory
func RemoteChromeDefault() (*RemoteConfig, context.CancelFunc, error) {
	rc := RemoteConfigDefault()
	p, err := rc.remoteChrome()
	return rc, p, err
}

func (rc *RemoteConfig) remoteChrome() (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(
		ctx,
		rc.ExecName,
		fmt.Sprintf("--remote-debugging-port=%d", rc.Port),
		fmt.Sprintf("--user-data-dir=%s", rc.UserDataDir),
		fmt.Sprintf("http://%s:%d%s", rc.OriginAddr, rc.OriginPort, rc.OriginRoute),
	)

	return cancel, cmd.Start()
}

// ReloadAllTabs will reload all opened tabs
func (rc *RemoteConfig) ReloadAllTabs() error {
	tabs, err := getTabs(rc.OriginAddr, rc.Port)
	if err != nil {
		return fmt.Errorf("error while getting tabs: %s", err)
	}
	var e error
	for _, tab := range tabs {
		err := reloadTab(tab)
		if err != nil {
			e = fmt.Errorf("error while reloading tab: %s", err)
		}
	}
	return e
}

// ReloadTab reloads one chrome tab by checking if the tab URL has route as suffix
func (rc *RemoteConfig) ReloadTab(route string) error {
	tabs, err := getTabs(rc.OriginAddr, rc.Port)
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
func (rc *RemoteConfig) ReloadTabGroup(subroute string) error {
	tabs, err := getTabs(rc.OriginAddr, rc.Port)
	if err != nil {
		return err
	}

	var e error
	for _, tab := range tabs {
		if strings.Contains(tab.URL, subroute) {
			err := reloadTab(tab)
			if err != nil {
				e = err
			}
		}
	}
	return e

}

func reloadTab(tab chromeTab) error {
	url := tab.WebSocketDebuggerURL
	origin := tab.URL
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return fmt.Errorf("error while connecting to socket: %s", err)
	}

	jsonStruct := &refreshJSON{ID: 0, Method: "Page.reload", Params: parameters{IgnoreCache: true}}
	jsonString, err := json.Marshal(jsonStruct)
	if err != nil {
		return fmt.Errorf("error while marshalling json: %s", err)
	}
	_, err = ws.Write(jsonString)
	if err != nil {
		return fmt.Errorf("error while preparing json: %s", err)
	}
	return nil
}

func getTabs(addr string, port int) ([]chromeTab, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/json", addr, port))
	if err != nil {
		return nil, fmt.Errorf("error while getting tabs: %s", err)
	}
	defer resp.Body.Close()

	var tabs []chromeTab
	buffer, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buffer, &tabs)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling tabs: %s", err)
	}
	return tabs, nil
}
