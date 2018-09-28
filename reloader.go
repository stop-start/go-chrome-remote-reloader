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
	// DebugPort is the chrome remote debugging port, usually 9222
	DebugPort int
	// UserDataDir is the folder for chrome user profile
	UserDataDir string
	// Addr is the base address to be opened
	Addr string
	// Port is the port for the base address to be opened
	Port int
	// Route path to be opened
	Route string
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

// NewRemoteConfig returns a default config for RemoteChrome
func NewRemoteConfig() *RemoteConfig {
	return &RemoteConfig{
		ExecName:    "google-chrome",
		DebugPort:   9222,
		UserDataDir: "/tmp/.chrome-remote-profile",
		Addr:        "localhost",
		Port:        8080,
		Route:       "",
	}
}

// RemoteChrome starts a new chrome remote debugging session
func (rc *RemoteConfig) RemoteChrome() (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(
		ctx,
		rc.ExecName,
		fmt.Sprintf("--remote-debugging-port=%d", rc.Port),
		fmt.Sprintf("--user-data-dir=%s", rc.UserDataDir),
		fmt.Sprintf("http://%s:%d%s", rc.Addr, rc.Port, rc.Route),
	)

	return cancel, cmd.Start()
}

// ReloadAllTabs will reload all opened tabs
func (rc *RemoteConfig) ReloadAllTabs() error {
	return reloadTabs(rc.Addr, rc.DebugPort, nil, "")
}

// ReloadTab reloads one chrome tab by checking if the tab URL has route as suffix
func (rc *RemoteConfig) ReloadTab(route string) error {
	return reloadTabs(rc.Addr, rc.DebugPort, strings.HasSuffix, route)
}

// ReloadTabGroup reloads a group of chrome tabs by checking if the tab URL contains the subroute
func (rc *RemoteConfig) ReloadTabGroup(subroute string) error {
	return reloadTabs(rc.Addr, rc.DebugPort, strings.Contains, subroute)
}

func reloadTabs(addr string, port int, conditionFunc func(string, string) bool, route string) error {
	tabs, err := getTabs(addr, port)
	if err != nil {
		return err
	}

	var e error
	for _, tab := range tabs {
		if conditionFunc == nil || conditionFunc(tab.URL, route) {
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
